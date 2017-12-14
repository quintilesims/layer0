package aws

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/zpatrick/go-bytesize"
	cache "github.com/zpatrick/go-cache"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/api/scaler"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type EnvironmentScaler struct {
	Client              *awsc.Client
	EnvironmentProvider provider.EnvironmentProvider
	ServiceProvider     provider.ServiceProvider
	TaskProvider        provider.TaskProvider
	JobStore            job.Store
	Config              config.APIConfig
	deployCache         *cache.Cache
	environmentCache    *cache.Cache
}

func NewEnvironmentScaler(a *awsc.Client, e provider.EnvironmentProvider, s provider.ServiceProvider, t provider.TaskProvider, j job.Store, c config.APIConfig) *EnvironmentScaler {
	return &EnvironmentScaler{
		Client:              a,
		EnvironmentProvider: e,
		ServiceProvider:     s,
		TaskProvider:        t,
		JobStore:            j,
		Config:              c,
		deployCache:         cache.New(),
		environmentCache:    cache.New(),
	}
}

type ProviderDistribution struct {
	Errs      []error
	Providers []*scaler.ResourceProvider
	SortBy    string
}

// Scale determines whether or not instances need to be added to or removed from a Layer0 environment, and makes any necessary changes.
// It consists of three primary logical groupings:
// 1. Gather all providers (instances) and consumers (tasks/services) of resources.
// 2. Calculate the optimal distribution of consumers among providers, including whether instances should be added or removed.
// 3. Update the AutoScaling Group to realize the changes calculated previously.
func (e *EnvironmentScaler) Scale(environmentID string) error {
	clusterName := addLayer0Prefix(e.Config.Instance(), environmentID)

	resourceConsumers, resourceProviders, err := e.GetCurrentState(clusterName)
	if err != nil {
		return err
	}

	resourceProviders, unusedProviders, desiredScale, calcErrs, err := e.CalculateOptimizedState(clusterName, resourceConsumers, resourceProviders)
	if err != nil {
		return err
	}

	if err := e.ScaleToState(clusterName, desiredScale, unusedProviders); err != nil {
		return err
	}

	if len(calcErrs) > 0 {
		return errors.MultiError(calcErrs)
	}

	return nil
}

func (e *EnvironmentScaler) GetCurrentState(clusterName string) ([]scaler.ResourceConsumer, []*scaler.ResourceProvider, error) {
	resourceProviders, err := e.getResourceProviders(clusterName)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("[DEBUG] [EnvironmentScaler] resourceProviders for env '%s': %v", clusterName, resourceProviders)

	var resourceConsumers []scaler.ResourceConsumer

	serviceConsumers, err := e.getServiceResourceConsumers(clusterName)
	if err != nil {
		return nil, nil, err
	}

	resourceConsumers = append(resourceConsumers, serviceConsumers...)

	ecsTaskConsumers, err := e.getECSTaskResourceConsumers(clusterName)
	if err != nil {
		return nil, nil, err
	}

	resourceConsumers = append(resourceConsumers, ecsTaskConsumers...)

	jobTaskConsumers, err := e.getJobTaskResourceConsumers(clusterName)
	if err != nil {
		return nil, nil, err
	}

	resourceConsumers = append(resourceConsumers, jobTaskConsumers...)

	log.Printf("[DEBUG] [EnvironmentScaler] resourceConsumers for env '%s': %#v", clusterName, resourceConsumers)

	return resourceConsumers, resourceProviders, nil
}

func (e *EnvironmentScaler) CalculateOptimizedState(clusterName string, resourceConsumers []scaler.ResourceConsumer, resourceProviders []*scaler.ResourceProvider) ([]*scaler.ResourceProvider, []*scaler.ResourceProvider, int, []error, error) {
	// calculate for scaling up
	resourceProviders, scaleUpErrs, err := e.calculateScaleUp(clusterName, resourceProviders, resourceConsumers)
	if err != nil {
		return nil, nil, 0, nil, err
	}

	// calculate for scaling down
	unusedProviders := e.calculateScaleDown(clusterName, resourceProviders)

	desiredScale := len(resourceProviders) - len(unusedProviders)

	return resourceProviders, unusedProviders, desiredScale, scaleUpErrs, nil
}

func (e *EnvironmentScaler) ScaleToState(clusterName string, desiredScale int, unusedProviders []*scaler.ResourceProvider) error {
	asg, err := readASG(e.Client.AutoScaling, clusterName)
	if err != nil {
		return err
	}

	currentCapacity := int(aws.Int64Value(asg.DesiredCapacity))
	if desiredScale == currentCapacity {
		log.Printf("[DEBUG] {EnvironmentScaler] Environment '%s' is at desired scale of '%d'. No scaling action required.", clusterName, currentCapacity)
		return nil
	}

	log.Printf("[DEBUG] [EnvironmentScaler] Attempting to scale environment '%s' from current scale of '%d' to desired scale of '%d'.", clusterName, currentCapacity, desiredScale)

	minCapacity := int(aws.Int64Value(asg.MinSize))
	if desiredScale < minCapacity {
		log.Printf("[DEBUG] [EnvironmentScaler] Will not scale below minimum capacity of '%d'. Aborting scaling action for environment '%s'.", minCapacity, clusterName)
		return nil
	}

	maxCapacity := int(aws.Int64Value(asg.MaxSize))
	if desiredScale > maxCapacity {
		log.Printf("[DEBUG] [EnvironmentScaler] Will not scale above maximum capacity of '%d'. Aborting scaling action for environment '%s'.", maxCapacity, clusterName)
		return nil
	}

	asgName := aws.StringValue(asg.AutoScalingGroupName)
	desiredScale64 := int64(desiredScale)
	if err := updateASG(e.Client.AutoScaling, asgName, nil, nil, &desiredScale64); err != nil {
		return err
	}

	// choose which instances to terminate during our scale-down process
	// instead of having asg randomly select instances
	// e.g. if we scale from 5 to 3, we can terminate up to 2 unused instances
	if currentCapacity > desiredScale {
		maxNumberOfInstancesToTerminate := currentCapacity - desiredScale

		canTerminate := func(i int) bool {
			if i+1 > maxNumberOfInstancesToTerminate {
				return false
			}

			if i > len(unusedProviders)-1 {
				return false
			}

			return true
		}

		for i := 0; canTerminate(i); i++ {
			unusedProvider := unusedProviders[i]
			instanceID := unusedProvider.ID

			log.Printf("[DEBUG] [EnvironmentScaler] Terminating unused instance '%s' from environment '%s'.", instanceID, clusterName)

			if err := e.terminateInstanceInAutoScalingGroup(instanceID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *EnvironmentScaler) calculateNewProvider(clusterName string) (*scaler.ResourceProvider, error) {
	var environment *models.Environment

	val, ok := e.environmentCache.Getf(clusterName)

	switch ok {
	case true:
		environment = val.(*models.Environment)

	case false:
		env, err := e.EnvironmentProvider.Read(clusterName)
		if err != nil {
			return nil, err
		}

		e.environmentCache.Add(clusterName, env)
		environment = env
	}

	instanceSpec, ok := instanceSpecifications()[environment.InstanceSize]
	if !ok {
		return nil, fmt.Errorf("[EnvironmentScaler] Instance size '%s' is not valid!", environment.InstanceSize)
	}

	return scaler.NewResourceProvider(instanceSpec.CPU, "<new instance>", instanceSpec.Memory), nil
}

func (e *EnvironmentScaler) calculateScaleDown(clusterName string, resourceProviders []*scaler.ResourceProvider) []*scaler.ResourceProvider {
	var unusedProviders []*scaler.ResourceProvider

	for _, provider := range resourceProviders {
		if !provider.InUse {
			unusedProviders = append(unusedProviders, provider)
		}
	}

	return unusedProviders
}

func (e *EnvironmentScaler) calculateScaleUp(clusterName string, resourceProviders []*scaler.ResourceProvider, resourceConsumers []scaler.ResourceConsumer) ([]*scaler.ResourceProvider, []error, error) {
	sorts := []ProviderDistribution{
		{[]error{}, resourceProviders, "cpu"},
		{[]error{}, resourceProviders, "mem"},
	}

	for _, s := range sorts {
		switch s.SortBy {
		case "cpu":
			sortProvidersByCPU(s.Providers)
			sortConsumersByCPU(resourceConsumers)

		case "mem":
			sortProvidersByMemory(s.Providers)
			sortConsumersByMemory(resourceConsumers)
		}

		sortProvidersByUsage(s.Providers)

		for _, consumer := range resourceConsumers {
			hasRoom := false

			for _, provider := range s.Providers {
				if provider.HasResourcesFor(consumer) {
					hasRoom = true
					if err := provider.SubtractResourcesFor(consumer); err != nil {
						s.Errs = append(s.Errs, err)
					}

					break
				}
			}

			if !hasRoom {
				newProvider, err := e.calculateNewProvider(clusterName)
				if err != nil {
					return nil, nil, err
				}

				if !newProvider.HasResourcesFor(consumer) {
					text := fmt.Sprintf("Resource '%s' cannot fit into an empty provider!", consumer.ID)
					text += "\nThe instance size in your environment is too small to run this resource."
					text += "\n Please increase the instance size for your environment."
					err := fmt.Errorf(text)
					s.Errs = append(s.Errs, err)
					continue
				}

				if err := newProvider.SubtractResourcesFor(consumer); err != nil {
					s.Errs = append(s.Errs, err)
				}

				s.Providers = append(s.Providers, newProvider)
			}
		}
	}

	p := findOptimalProviderDistribution(sorts)

	return p.Providers, p.Errs, nil
}

func findOptimalProviderDistribution(providerDistributions []ProviderDistribution) ProviderDistribution {
	if len(providerDistributions) == 1 {
		return providerDistributions[0]
	}

	shortestLength := len(providerDistributions[0].Providers)
	shortestDistribution := providerDistributions[0]

	for _, p := range providerDistributions[1:] {
		l := len(p.Providers)
		if l < shortestLength {
			shortestLength = l
			shortestDistribution = p
		}
	}

	return shortestDistribution
}

func (e *EnvironmentScaler) getContainerResourceFromDeploy(deployID string) ([]scaler.ResourceConsumer, error) {
	if consumers, ok := e.deployCache.Getf(deployID); ok {
		return consumers.([]scaler.ResourceConsumer), nil
	}

	taskDefinition, err := e.getTaskDefinitionFromDeployID(deployID)
	if err != nil {
		return nil, err
	}

	consumers := make([]scaler.ResourceConsumer, len(taskDefinition.ContainerDefinitions))
	for i, taskDefinition := range taskDefinition.ContainerDefinitions {
		var cpu int
		var memory bytesize.Bytesize

		if c := int(aws.Int64Value(taskDefinition.Cpu)); c != 0 {
			cpu = c
		}

		if m := aws.Int64Value(taskDefinition.MemoryReservation); m != 0 {
			memory = bytesize.MiB * bytesize.Bytesize(m)
		}

		if m := aws.Int64Value(taskDefinition.Memory); m != 0 {
			memory = bytesize.MiB * bytesize.Bytesize(m)
		}

		ports := []int{}
		for _, p := range taskDefinition.PortMappings {
			if hostPort := int(aws.Int64Value(p.HostPort)); hostPort != 0 {
				ports = append(ports, hostPort)
			}
		}

		id := "NewConsumerFromDeploy" + deployID
		consumers[i] = scaler.NewResourceConsumer(cpu, id, memory, ports)
	}

	e.deployCache.Add(deployID, consumers)
	return consumers, nil
}

func (e *EnvironmentScaler) getServiceResourceConsumers(clusterName string) ([]scaler.ResourceConsumer, error) {
	var resourceConsumers []scaler.ResourceConsumer

	serviceARNs, err := e.getServiceARNsForCluster(clusterName)
	if err != nil {
		return nil, err
	}

	services, err := e.getServicesFromServiceARNs(clusterName, serviceARNs)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		deployIDCopies := map[string]int64{}
		// deployment.RunningCount is the number of containers already running on an instance
		// deployment.PendingCount is the number of containers that are alraedy on an instance, but are being pulled
		// we only care about containers that are not on instances yet
		for _, d := range service.Deployments {
			desiredCount := aws.Int64Value(d.DesiredCount)
			runningCount := aws.Int64Value(d.RunningCount)
			pendingCount := aws.Int64Value(d.PendingCount)
			if numPending := desiredCount - (runningCount + pendingCount); numPending > 0 {
				deployIDCopies[aws.StringValue(d.Id)] = numPending
			}
		}

		if len(deployIDCopies) == 0 {
			continue
		}

		// iterate through deploys
		for deployID := range deployIDCopies {
			c, err := e.getContainerResourceFromDeploy(deployID)
			if err != nil {
				return nil, err
			}

			resourceConsumers = append(resourceConsumers, c...)
		}
	}

	return resourceConsumers, nil
}

func (e *EnvironmentScaler) getECSTaskResourceConsumers(clusterName string) ([]scaler.ResourceConsumer, error) {
	var taskARNs []string

	startedBy := e.Config.Instance()
	for _, status := range []string{ecs.DesiredStatusRunning, ecs.DesiredStatusStopped} {
		arns, err := e.getTaskARNsForCluster(clusterName, status, startedBy)
		if err != nil {
			return nil, err
		}

		taskARNs = append(taskARNs, arns...)
	}

	var resourceConsumers []scaler.ResourceConsumer

	for _, taskARN := range taskARNs {
		task, err := e.TaskProvider.Read(taskARN)
		if err != nil {
			return nil, err
		}

		taskResourceConsumers, err := e.getContainerResourceFromDeploy(task.DeployID)
		if err != nil {
			return nil, err
		}

		resourceConsumers = append(resourceConsumers, taskResourceConsumers...)
	}

	return resourceConsumers, nil
}

func (e *EnvironmentScaler) getJobTaskResourceConsumers(clusterName string) ([]scaler.ResourceConsumer, error) {
	var resourceConsumers []scaler.ResourceConsumer
	jobs, err := e.JobStore.SelectAll()
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		if job.Type == models.CreateTaskJob {
			// Don't check for Pending jobs.
			// It's more efficient and cost-effective to only scale for jobs that are InProgress.
			// Jobs could be stuck in the Pending state for hours waiting for a worker to pick them up.
			if job.Status == models.InProgressJobStatus {
				var req models.CreateTaskRequest
				if err := json.Unmarshal([]byte(job.Request), &req); err != nil {
					return nil, err
				}

				if req.EnvironmentID == clusterName {
					taskResourceConsumers, err := e.getContainerResourceFromDeploy(req.DeployID)
					if err != nil {
						return nil, err
					}

					resourceConsumers = append(resourceConsumers, taskResourceConsumers...)
				}
			}
		}
	}

	return resourceConsumers, nil
}

func (e *EnvironmentScaler) getResourceProviders(clusterName string) ([]*scaler.ResourceProvider, error) {
	var result []*scaler.ResourceProvider

	containerInstanceARNs, err := e.getActiveContainerInstanceARNsForCluster(clusterName)
	if err != nil {
		return nil, err
	}

	if len(containerInstanceARNs) == 0 {
		return result, nil
	}

	containerInstances, err := e.getContainerInstancesFromARNs(clusterName, containerInstanceARNs)
	if err != nil {
		return nil, err
	}

	if len(containerInstances) == 0 {
		return result, nil
	}

	for _, instance := range containerInstances {
		// it's non-intuitive, but the ports being used by the tasks live in
		// instance.RemainingResources, not instance.RegisteredResources
		var usedPorts []int
		var availableCPU int
		var availableMemory bytesize.Bytesize

		for _, resource := range instance.RemainingResources {
			switch aws.StringValue(resource.Name) {
			case "CPU":
				availableCPU = int(aws.Int64Value(resource.IntegerValue))

			case "MEMORY":
				val := aws.Int64Value(resource.IntegerValue)
				availableMemory = bytesize.MiB * bytesize.Bytesize(val)

			case "PORTS":
				for _, port := range resource.StringSetValue {
					port, err := strconv.Atoi(aws.StringValue(port))
					if err != nil {
						log.Printf("[WARN] [EnvironmentScaler] Instance %s: Failed to convert port to int: %v", aws.StringValue(instance.Ec2InstanceId), err)
						continue
					}

					usedPorts = append(usedPorts, port)
				}
			}
		}

		inUse := aws.Int64Value(instance.PendingTasksCount)+aws.Int64Value(instance.RunningTasksCount) > 0

		r := &scaler.ResourceProvider{
			AgentIsConnected: aws.BoolValue(instance.AgentConnected),
			AvailableCPU:     availableCPU,
			AvailableMemory:  availableMemory,
			ID:               aws.StringValue(instance.Ec2InstanceId),
			InUse:            inUse,
			Status:           aws.StringValue(instance.Status),
			UsedPorts:        usedPorts,
		}

		result = append(result, r)
	}

	return result, nil
}

func sortConsumersByCPU(c []scaler.ResourceConsumer) {
	sorter := &ResourceConsumerSorter{
		Consumers: c,
		lessThan: func(i scaler.ResourceConsumer, j scaler.ResourceConsumer) bool {
			return i.CPU < j.CPU
		},
	}

	sort.Sort(sorter)
}

func sortConsumersByMemory(c []scaler.ResourceConsumer) {
	sorter := &ResourceConsumerSorter{
		Consumers: c,
		lessThan: func(i scaler.ResourceConsumer, j scaler.ResourceConsumer) bool {
			return i.Memory < j.Memory
		},
	}

	sort.Sort(sorter)
}

func sortProvidersByCPU(p []*scaler.ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: p,
		lessThan: func(i *scaler.ResourceProvider, j *scaler.ResourceProvider) bool {
			return i.AvailableCPU < j.AvailableCPU
		},
	}

	sort.Sort(sorter)
}

func sortProvidersByMemory(p []*scaler.ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: p,
		lessThan: func(i *scaler.ResourceProvider, j *scaler.ResourceProvider) bool {
			return i.AvailableMemory < j.AvailableMemory
		},
	}

	sort.Sort(sorter)
}

func sortProvidersByUsage(r []*scaler.ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: r,
		lessThan: func(i *scaler.ResourceProvider, j *scaler.ResourceProvider) bool {
			return i.InUse && !j.InUse
		},
	}

	sort.Sort(sorter)
}

type ResourceConsumerSorter struct {
	Consumers []scaler.ResourceConsumer
	lessThan  func(scaler.ResourceConsumer, scaler.ResourceConsumer) bool
}

func (r *ResourceConsumerSorter) Len() int {
	return len(r.Consumers)
}

func (r *ResourceConsumerSorter) Swap(i, j int) {
	r.Consumers[i], r.Consumers[j] = r.Consumers[j], r.Consumers[i]
}

func (r *ResourceConsumerSorter) Less(i, j int) bool {
	return r.lessThan(r.Consumers[i], r.Consumers[j])
}

type ResourceProviderSorter struct {
	Providers []*scaler.ResourceProvider
	lessThan  func(*scaler.ResourceProvider, *scaler.ResourceProvider) bool
}

func (r *ResourceProviderSorter) Len() int {
	return len(r.Providers)
}

func (r *ResourceProviderSorter) Swap(i, j int) {
	r.Providers[i], r.Providers[j] = r.Providers[j], r.Providers[i]
}

func (r *ResourceProviderSorter) Less(i, j int) bool {
	return r.lessThan(r.Providers[i], r.Providers[j])
}
