package aws

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/zpatrick/go-bytesize"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/api/scaler"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func defaultPorts() []int {
	return []int{
		22,
		2376,
		2375,
		51678,
		51679,
	}
}

type InstanceSpec struct {
	CPU    int
	Memory bytesize.Bytesize
}

func NewInstanceSpec(virtualCPUCores int, gibibytesMemory float64) InstanceSpec {
	return InstanceSpec{
		virtualCPUCores,
		bytesize.GiB * bytesize.Bytesize(gibibytesMemory),
	}
}

// https://aws.amazon.com/ec2/instance-types/
func instanceSpecs() map[string]InstanceSpec {
	return map[string]InstanceSpec{
		"t2.nano":      NewInstanceSpec(1, 0.5),
		"t2.micro":     NewInstanceSpec(1, 1),
		"t2.small":     NewInstanceSpec(1, 2),
		"t2.medium":    NewInstanceSpec(2, 4),
		"t2.large":     NewInstanceSpec(2, 8),
		"t2.xlarge":    NewInstanceSpec(4, 16),
		"t2.2xlarge":   NewInstanceSpec(8, 32),
		"m5.large":     NewInstanceSpec(2, 8),
		"m5.xlarge":    NewInstanceSpec(4, 16),
		"m5.2xlarge":   NewInstanceSpec(8, 32),
		"m5.4xlarge":   NewInstanceSpec(16, 64),
		"m5.12xlarge":  NewInstanceSpec(48, 192),
		"m5.24xlarge":  NewInstanceSpec(96, 384),
		"m4.large":     NewInstanceSpec(2, 8),
		"m4.xlarge":    NewInstanceSpec(4, 16),
		"m4.2xlarge":   NewInstanceSpec(8, 32),
		"m4.4xlarge":   NewInstanceSpec(16, 64),
		"m4.10xlarge":  NewInstanceSpec(40, 160),
		"m4.16xlarge":  NewInstanceSpec(64, 256),
		"m3.medium":    NewInstanceSpec(1, 3.75),
		"m3.large":     NewInstanceSpec(2, 7.5),
		"m3.xlarge":    NewInstanceSpec(4, 15),
		"m3.2xlarge":   NewInstanceSpec(8, 30),
		"c5.large":     NewInstanceSpec(2, 4),
		"c5.xlarge":    NewInstanceSpec(4, 8),
		"c5.2xlarge":   NewInstanceSpec(8, 16),
		"c5.4xlarge":   NewInstanceSpec(16, 32),
		"c5.9xlarge":   NewInstanceSpec(36, 72),
		"c5.18xlarge":  NewInstanceSpec(72, 144),
		"c4.large":     NewInstanceSpec(2, 3.75),
		"c4.xlarge":    NewInstanceSpec(4, 7.5),
		"c4.2xlarge":   NewInstanceSpec(8, 15),
		"c4.4xlarge":   NewInstanceSpec(16, 30),
		"c4.8xlarge":   NewInstanceSpec(36, 60),
		"c3.large":     NewInstanceSpec(2, 3.75),
		"c3.xlarge":    NewInstanceSpec(4, 7.5),
		"c3.2xlarge":   NewInstanceSpec(8, 15),
		"c3.4xlarge":   NewInstanceSpec(16, 30),
		"c3.8xlarge":   NewInstanceSpec(32, 60),
		"x1.16large":   NewInstanceSpec(64, 976),
		"x1.32xlarge":  NewInstanceSpec(128, 1952),
		"x1e.xlarge":   NewInstanceSpec(4, 122),
		"x1e.2xlarge":  NewInstanceSpec(8, 244),
		"x1e.4xlarge":  NewInstanceSpec(16, 488),
		"x1e.8xlarge":  NewInstanceSpec(32, 976),
		"x1e.16xlarge": NewInstanceSpec(64, 1952),
		"x1e.32xlarge": NewInstanceSpec(128, 3904),
		"r4.large":     NewInstanceSpec(2, 15.25),
		"r4.xlarge":    NewInstanceSpec(4, 30.5),
		"r4.2xlarge":   NewInstanceSpec(8, 61),
		"r4.4xlarge":   NewInstanceSpec(16, 122),
		"r4.8xlarge":   NewInstanceSpec(32, 244),
		"r4.16xlarge":  NewInstanceSpec(64, 488),
		"r3.large":     NewInstanceSpec(2, 15.25),
		"r3.xlarge":    NewInstanceSpec(4, 30.5),
		"r3.2xlarge":   NewInstanceSpec(8, 61),
		"r3.4xlarge":   NewInstanceSpec(16, 122),
		"r3.8xlarge":   NewInstanceSpec(32, 244),
		"p3.2xlarge":   NewInstanceSpec(8, 61),
		"p3.8xlarge":   NewInstanceSpec(32, 244),
		"p3.16xlarge":  NewInstanceSpec(64, 488),
		"p2.xlarge":    NewInstanceSpec(4, 61),
		"p2.8xlarge":   NewInstanceSpec(32, 488),
		"p2.16xlarge":  NewInstanceSpec(64, 732),
		"g3.4xlarge":   NewInstanceSpec(16, 122),
		"g3.8xlarge":   NewInstanceSpec(32, 244),
		"g3.16xlarge":  NewInstanceSpec(64, 488),
		"f1.2xlarge":   NewInstanceSpec(8, 122),
		"f1.16xlarge":  NewInstanceSpec(64, 976),
		"h1.2xlarge":   NewInstanceSpec(8, 32),
		"h1.4xlarge":   NewInstanceSpec(16, 64),
		"h1.8xlarge":   NewInstanceSpec(32, 128),
		"h1.16xlarge":  NewInstanceSpec(64, 256),
		"i3.large":     NewInstanceSpec(2, 15.25),
		"i3.xlarge":    NewInstanceSpec(4, 30.5),
		"i3.2xlarge":   NewInstanceSpec(8, 61),
		"i3.4xlarge":   NewInstanceSpec(16, 122),
		"i3.8xlarge":   NewInstanceSpec(32, 244),
		"i3.16xlarge":  NewInstanceSpec(64, 488),
		"i3.metal":     NewInstanceSpec(0, 512),
		"d2.xlarge":    NewInstanceSpec(4, 30.5),
		"d2.2xlarge":   NewInstanceSpec(8, 61),
		"d2.4xlarge":   NewInstanceSpec(16, 122),
		"d2.8xlarge":   NewInstanceSpec(36, 244),
	}
}

type EnvironmentScaler struct {
	Client              *awsc.Client
	EnvironmentProvider provider.EnvironmentProvider
	ServiceProvider     provider.ServiceProvider
	TaskProvider        provider.TaskProvider
	JobStore            job.Store
	Config              config.APIConfig
	deployCache         map[string][]scaler.ResourceConsumer
}

func NewEnvironmentScaler(a *awsc.Client, e provider.EnvironmentProvider, s provider.ServiceProvider, t provider.TaskProvider, j job.Store, c config.APIConfig) *EnvironmentScaler {
	return &EnvironmentScaler{
		Client:              a,
		EnvironmentProvider: e,
		ServiceProvider:     s,
		TaskProvider:        t,
		JobStore:            j,
		Config:              c,
		// deployCache: map[string][]*Resource,
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

	resourceProviders, unusedProviders, calcErrs, err := e.CalculateOptimizedState(clusterName, resourceConsumers, resourceProviders)
	if err != nil {
		return err
	}

	desiredScale := len(resourceProviders) - len(unusedProviders)
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

	humanReadableProviders := ""
	for _, r := range resourceProviders {
		humanReadableProviders += fmt.Sprintf("\n    %#v", r)
	}

	log.Printf("[DEBUG] [EnvironmentScaler] resourceProviders for env '%s': %s", clusterName, humanReadableProviders)

	var resourceConsumers []scaler.ResourceConsumer

	serviceConsumers, err := e.getResourceConsumers_PendingServices(clusterName)
	if err != nil {
		return nil, nil, err
	}

	resourceConsumers = append(resourceConsumers, serviceConsumers...)

	ecsTaskConsumers, err := e.getResourceConsumers_TasksInECS(clusterName)
	if err != nil {
		return nil, nil, err
	}

	resourceConsumers = append(resourceConsumers, ecsTaskConsumers...)

	jobTaskConsumers, err := e.getResourceConsumers_TasksInJobs(clusterName)
	if err != nil {
		return nil, nil, err
	}

	resourceConsumers = append(resourceConsumers, jobTaskConsumers...)

	log.Printf("[DEBUG] [EnvironmentScaler] resourceConsumers for env '%s': %#v", clusterName, resourceConsumers)

	return resourceConsumers, resourceProviders, nil
}

func (e *EnvironmentScaler) CalculateOptimizedState(clusterName string, resourceConsumers []scaler.ResourceConsumer, resourceProviders []*scaler.ResourceProvider) ([]*scaler.ResourceProvider, []*scaler.ResourceProvider, []error, error) {
	// calculate for scaling up
	resourceProviders, scaleUpErrs, err := e.calculateScaleUp(clusterName, resourceProviders, resourceConsumers)
	if err != nil {
		return nil, nil, nil, err
	}

	// calculate for scaling down
	unusedProviders := e.calculateScaleDown(clusterName, resourceProviders)

	return resourceProviders, unusedProviders, scaleUpErrs, nil
}

func (e *EnvironmentScaler) ScaleToState(clusterName string, desiredScale int, unusedProviders []*scaler.ResourceProvider) error {
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	input.SetAutoScalingGroupNames([]*string{&clusterName})

	asgs, err := e.Client.AutoScaling.DescribeAutoScalingGroups(input)
	if err != nil {
		return err
	}

	asg := asgs.AutoScalingGroups[0]

	currentCapacity := int(aws.Int64Value(asg.DesiredCapacity))

	switch {
	case desiredScale > currentCapacity:
		log.Printf("[DEBUG] [EnvironmentScaler] Attempting to scale environment '%s' from current scale of '%d' to desired scale of '%d'.", clusterName, currentCapacity, desiredScale)
		return e.scaleUp(clusterName, desiredScale, asg)

	case desiredScale < currentCapacity:
		log.Printf("[DEBUG] [EnvironmentScaler] Attempting to scale environment '%s' from current scale of '%d' to desired scale of '%d'.", clusterName, currentCapacity, desiredScale)
		return e.scaleDown(clusterName, desiredScale, asg, unusedProviders)

	default:
		log.Printf("[DEBUG] [EnvironmentScaler] Environment '%s' is at desired scale of '%d'. No scaling action required.", clusterName, currentCapacity)
		return nil
	}
}

func (e *EnvironmentScaler) calculateNewProvider(clusterName string) (*scaler.ResourceProvider, error) {
	env, err := e.EnvironmentProvider.Read(clusterName)
	if err != nil {
		return nil, err
	}

	instanceSpec := instanceSpecs()[env.InstanceSize]

	resource := &scaler.ResourceProvider{}
	resource.AvailableCPU = instanceSpec.CPU
	resource.AvailableMemory = instanceSpec.Memory
	resource.ID = "<new instance>"
	resource.InUse = false
	resource.UsedPorts = defaultPorts()

	return resource, nil
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
					provider.SubtractResourcesFor(consumer)
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

				newProvider.SubtractResourcesFor(consumer)
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
	// use some kind of deploy cache
	if consumers, ok := e.deployCache[deployID]; ok {
		return consumers, nil
	}

	input := &ecs.DescribeTaskDefinitionInput{}
	input.SetTaskDefinition(deployID)
	output, err := e.Client.ECS.DescribeTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	consumers := make([]scaler.ResourceConsumer, len(output.TaskDefinition.ContainerDefinitions))
	for i, d := range output.TaskDefinition.ContainerDefinitions {
		var memory bytesize.Bytesize

		if aws.Int64Value(d.MemoryReservation) != 0 {
			memory = bytesize.MiB * bytesize.Bytesize(aws.Int64Value(d.MemoryReservation))
		}

		if aws.Int64Value(d.Memory) != 0 {
			memory = bytesize.MiB * bytesize.Bytesize(aws.Int64Value(d.Memory))
		}

		ports := []int{}
		for _, p := range d.PortMappings {
			if aws.Int64Value(p.HostPort) != 0 {
				ports = append(ports, int(aws.Int64Value(p.HostPort)))
			}
		}

		consumers[i] = scaler.ResourceConsumer{
			ID:     "",
			Memory: memory,
			Ports:  ports,
		}
	}

	e.deployCache[deployID] = consumers
	return consumers, nil
}

func (e *EnvironmentScaler) getResourceConsumers_PendingServices(clusterName string) ([]scaler.ResourceConsumer, error) {
	var resourceConsumers []scaler.ResourceConsumer

	listServicesInput := &ecs.ListServicesInput{}
	listServicesInput.SetCluster(clusterName)

	serviceARNs := []*string{}
	listServicesPagesFN := func(output *ecs.ListServicesOutput, lastPage bool) bool {
		serviceARNs = append(serviceARNs, output.ServiceArns...)

		return !lastPage
	}

	if err := e.Client.ECS.ListServicesPages(listServicesInput, listServicesPagesFN); err != nil {
		return nil, err
	}

	services := []*ecs.Service{}
	if len(serviceARNs) > 0 {
		// The SDK states that you can specify up to 10 services in one DescribeServices operation:
		// https://github.com/aws/aws-sdk-go/blob/v1.12.19/service/ecs/api.go#L5420
		// (aws-sdk-go version 1.12.19, as stated in layer0/Gopkg.toml)
		for i := 0; i < len(serviceARNs); i += 10 {
			end := i + 10

			if end > len(serviceARNs) {
				end = len(serviceARNs)
			}

			describeServicesInput := &ecs.DescribeServicesInput{}
			describeServicesInput.SetCluster(clusterName)
			describeServicesInput.SetServices(serviceARNs[i:end])

			output, err := e.Client.ECS.DescribeServices(describeServicesInput)
			if err != nil {
				return nil, err
			}

			services = append(services, output.Services...)
		}
	}

	for _, service := range services {
		deployIDCopies := map[string]int64{}
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

func (e *EnvironmentScaler) getResourceConsumers_TasksInECS(clusterName string) ([]scaler.ResourceConsumer, error) {
	var (
		taskARNs          []string
		resourceConsumers []scaler.ResourceConsumer
	)

	fn := func(output *ecs.ListTasksOutput, lastPage bool) bool {

		for _, taskARN := range output.TaskArns {
			taskARNs = append(taskARNs, aws.StringValue(taskARN))
		}

		return !lastPage
	}

	startedBy := e.Config.Instance()
	for _, status := range []string{ecs.DesiredStatusRunning, ecs.DesiredStatusStopped} {
		input := &ecs.ListTasksInput{}
		input.SetCluster(clusterName)
		input.SetDesiredStatus(status)
		input.SetStartedBy(startedBy)
		if err := e.Client.ECS.ListTasksPages(input, fn); err != nil {
			return nil, err
		}
	}

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

func (e *EnvironmentScaler) getResourceConsumers_TasksInJobs(clusterName string) ([]scaler.ResourceConsumer, error) {
	var resourceConsumers []scaler.ResourceConsumer
	jobs, err := e.JobStore.SelectAll()
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		if job.Type == models.CreateTaskJob {
			// don't check for Pending jobs; once the job runner has picked
			// up a job, its status is already InProgress
			if job.Status == models.InProgressJobStatus {
				var req models.CreateTaskRequest
				if err := json.Unmarshal([]byte(job.Request), &req); err != nil {
					return nil, err
				}

				if req.EnvironmentID == clusterName {
					// note that this isn't exact if the job has started some, but not all of the tasks
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

	listContainerInstancesInput := &ecs.ListContainerInstancesInput{}
	listContainerInstancesInput.SetCluster(clusterName)
	listContainerInstancesInput.SetStatus(ecs.ContainerInstanceStatusActive)

	containerInstanceARNs := []*string{}
	listContainerInstancesPagesFN := func(output *ecs.ListContainerInstancesOutput, lastPage bool) bool {
		containerInstanceARNs = append(containerInstanceARNs, output.ContainerInstanceArns...)

		return !lastPage
	}

	if err := e.Client.ECS.ListContainerInstancesPages(listContainerInstancesInput, listContainerInstancesPagesFN); err != nil {
		return nil, err
	}

	if len(containerInstanceARNs) == 0 {
		return result, nil
	}

	describeContainerInstancesInput := &ecs.DescribeContainerInstancesInput{}
	describeContainerInstancesInput.SetCluster(clusterName)
	describeContainerInstancesInput.SetContainerInstances(containerInstanceARNs)

	output, err := e.Client.ECS.DescribeContainerInstances(describeContainerInstancesInput)
	if err != nil {
		return nil, err
	}

	if len(output.ContainerInstances) == 0 {
		return result, nil
	}

	for _, instance := range output.ContainerInstances {
		// it's non-intuitive, but the ports being used by the tasks live in
		// instance.RemainingResources, not instance.RegisteredResources
		var (
			usedPorts       []int
			availableCPU    int
			availableMemory bytesize.Bytesize
		)

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

func (e *EnvironmentScaler) scaleDown(clusterName string, desiredScale int, asg *autoscaling.Group, unusedProviders []*scaler.ResourceProvider) error {
	minCapacity := int(aws.Int64Value(asg.MinSize))
	if desiredScale < minCapacity {
		log.Printf("[DEBUG] [EnvironmentScaler] Will not scale below minimum capacity of '%d'. Aborting scaling action for environment '%s'.", minCapacity, clusterName)
		return nil
	}

	input := &autoscaling.SetDesiredCapacityInput{}
	input.SetAutoScalingGroupName(aws.StringValue(asg.AutoScalingGroupName))
	input.SetDesiredCapacity(int64(desiredScale))

	if _, err := e.Client.AutoScaling.SetDesiredCapacity(input); err != nil {
		return err
	}

	// choose which instances to terminate during our scale-down process
	// instead of having asg randomly select instances
	// e.g. if we scale from 5 to 3, we can terminate up to 2 unused instances
	currentCapacity := int(aws.Int64Value(asg.DesiredCapacity))
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
		log.Printf("[DEBUG] [EnvironmentScaler] Terminating unused instance '%s' from environment '%s'.", unusedProvider.ID, clusterName)

		input := &autoscaling.TerminateInstanceInAutoScalingGroupInput{}
		input.SetInstanceId(unusedProvider.ID)
		input.SetShouldDecrementDesiredCapacity(false)

		if _, err := e.Client.AutoScaling.TerminateInstanceInAutoScalingGroup(input); err != nil {
			return err
		}
	}

	return nil
}

func (e *EnvironmentScaler) scaleUp(clusterName string, desiredScale int, asg *autoscaling.Group) error {
	maxCapacity := int(aws.Int64Value(asg.MaxSize))
	if desiredScale > maxCapacity {
		log.Printf("[DEBUG] [EnvironmentScaler] Will not scale above maximum capacity of '%d'. Aborting scaling action for environment '%s'.", maxCapacity, clusterName)
		return nil
	}

	input := &autoscaling.SetDesiredCapacityInput{}
	input.SetAutoScalingGroupName(aws.StringValue(asg.AutoScalingGroupName))
	input.SetDesiredCapacity(int64(desiredScale))

	if _, err := e.Client.AutoScaling.SetDesiredCapacity(input); err != nil {
		return err
	}

	return nil
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
