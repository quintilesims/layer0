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
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

var defaultPorts = []int{
	22,
	2376,
	2375,
	51678,
	51679,
}

type EnvironmentScaler struct {
	Client              *awsc.Client
	EnvironmentProvider provider.EnvironmentProvider
	ServiceProvider     provider.ServiceProvider
	TaskProvider        provider.TaskProvider
	JobStore            job.Store
	Config              config.APIConfig
	deployCache         map[string][]*models.Resource
}

const (
	Memory = iota
	CPU
)

func NewEnvironmentScaler(a *awsc.Client, e provider.EnvironmentProvider, s provider.ServiceProvider, t provider.TaskProvider, j job.Store, c config.APIConfig) *EnvironmentScaler {
	return &EnvironmentScaler{
		Client:              a,
		EnvironmentProvider: e,
		ServiceProvider:     s,
		TaskProvider:        t,
		JobStore:            j,
		Config:              c,
		// deployCache: map[string][]*models.Resource,
	}
}

func (e *EnvironmentScaler) Scale(environmentID string) error {
	// GET RESOURCE PROVIDERS
	clusterName := addLayer0Prefix(e.Config.Instance(), environmentID)

	resourceProviders, err := e.getResourceProviders(clusterName)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] resourceProviders for env '%s': %#v, number of providers:%d", environmentID, resourceProviders, len(resourceProviders))

	// GET RESOURCE CONSUMERS
	resourceConsumers, err := e.getResourceConsumers(clusterName)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] resourceConsumers for env '%s': %#v, number of consumers:%d", environmentID, resourceConsumers, len(resourceConsumers))

	_, err = e.scale(clusterName, resourceProviders, resourceConsumers)
	if err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentScaler) scale(clusterName string, providers []*models.Resource, consumers []*models.Resource) (*models.ScalerRunInfo, error) {
	var errs []error

	// TODO: pick most efficient distribution of comsumers among providers
	// depending on whether we sort by memory or by CPU

	// check if we need to scale up
	for _, consumer := range consumers {
		hasRoom := false

		// first, sort by memory so we pack tasks by memory as tightly as possible
		sortProvidersByMemory(providers)

		// next, place any unused providers in the back of the list
		// that way, we can can delete them if we avoid placing any tasks in them
		sortProvidersByUsage(providers)

		for _, provider := range providers {

			if hasResourcesFor(consumer, provider) {
				hasRoom = true
				subtractResourcesFor(consumer, provider)
				break
			}
		}

		if !hasRoom {
			newProvider, err := e.calculateNewProvider(clusterName)
			if err != nil {
				return nil, err
			}

			if !hasResourcesFor(consumer, newProvider) {
				text := fmt.Sprintf("Resource '%s' cannot fit into an empty provider!", consumer.ID)
				text += "\nThe instance size in your environment is too small to run this resource."
				text += "\nPlease increase the instance size for your environment"
				err := fmt.Errorf(text)
				errs = append(errs, err)
				continue
			}

			subtractResourcesFor(consumer, newProvider)
			providers = append(providers, newProvider)
		}
	}

	// check if we need to scale down
	unusedProviders := []*models.Resource{}
	for _, provider := range providers {
		if !provider.InUse {
			unusedProviders = append(unusedProviders, provider)
		}
	}

	desiredScale := len(providers) - len(unusedProviders)

	actualScale, err := e.scaleTo(clusterName, desiredScale, unusedProviders)
	if err != nil {
		errs = append(errs, err)
	}

	scaleBeforeRun := len(providers)
	info := &models.ScalerRunInfo{
		EnvironmentID:           clusterName,
		PendingResources:        resourceModels(consumers),
		ResourceProviders:       resourceModels(providers),
		ScaleBeforeRun:          scaleBeforeRun,
		DesiredScaleAfterRun:    desiredScale,
		ActualScaleAfterRun:     actualScale,
		UnusedResourceProviders: len(unusedProviders),
	}

	return info, errors.MultiError(errs)
}

func resourceModels(resources []*models.Resource) []models.Resource {
	resourceModels := make([]models.Resource, len(resources))
	for i, resource := range resources {
		new := models.Resource{
			ID:     resource.ID,
			InUse:  resource.InUse,
			Ports:  resource.Ports,
			Memory: resource.Memory,
			CPU:    resource.CPU,
		}
		resourceModels[i] = new
	}

	return resourceModels
}

func (e *EnvironmentScaler) getResourceProviders(clusterName string) ([]*models.Resource, error) {
	listContainerInstancesInput := &ecs.ListContainerInstancesInput{}
	listContainerInstancesInput.SetCluster(clusterName)
	listContainerInstancesInput.SetStatus("ACTIVE")

	containerInstanceARNs := []*string{}
	listContainerInstancesPagesFN := func(output *ecs.ListContainerInstancesOutput, lastPage bool) bool {
		containerInstanceARNs = append(containerInstanceARNs, output.ContainerInstanceArns...)

		return !lastPage
	}

	if err := e.Client.ECS.ListContainerInstancesPages(listContainerInstancesInput, listContainerInstancesPagesFN); err != nil {
		return nil, err
	}

	describeContainerInstancesInput := &ecs.DescribeContainerInstancesInput{}
	describeContainerInstancesInput.SetCluster(clusterName)
	describeContainerInstancesInput.SetContainerInstances(containerInstanceARNs)

	output, err := e.Client.ECS.DescribeContainerInstances(describeContainerInstancesInput)
	if err != nil {
		return nil, err
	}

	result := []*models.Resource{}
	if len(containerInstanceARNs) == 0 {
		return result, nil
	}

	if len(output.ContainerInstances) == 0 {
		return result, nil
	}

	for _, instance := range output.ContainerInstances {
		if !aws.BoolValue(instance.AgentConnected) {
			log.Printf("[DEBUG] Not counting instance '%s' as a resource provider (ecs agent not connected)", aws.StringValue(instance.Ec2InstanceId))
			continue
		}

		if aws.StringValue(instance.Status) != "ACTIVE" {
			log.Printf("[DEBUG] Not counting instance '%s' as a resource provider (status != ACTIVE)", aws.StringValue(instance.Ec2InstanceId))
			continue
		}

		// it's non-intuitive, but the ports being used by the tasks live in
		// instance.RemainingResources, not instance.RegisteredResources
		var (
			usedPorts       []int
			availableMemory bytesize.Bytesize
			availableCPU    bytesize.Bytesize
		)

		for _, resource := range instance.RemainingResources {
			log.Printf("[DEBUG] instance.RemainingResources resource: %#v", resource)
			switch aws.StringValue(resource.Name) {
			// TODO: add CPU bound ("fun")
			case "MEMORY":
				val := aws.Int64Value(resource.IntegerValue)
				availableMemory = bytesize.MiB * bytesize.Bytesize(val)

			case "PORTS":
				for _, port := range resource.StringSetValue {
					port, err := strconv.Atoi(aws.StringValue(port))
					if err != nil {
						log.Printf("[WARN] Instance %s: Failed to convert port to int: %v", aws.StringValue(instance.Ec2InstanceId), err)
						continue
					}

					usedPorts = append(usedPorts, port)
				}
			case "CPU":
				val := aws.Int64Value(resource.IntegerValue)
				availableCPU = bytesize.MiB * bytesize.Bytesize(val)
			}
		}

		inUse := aws.Int64Value(instance.PendingTasksCount)+aws.Int64Value(instance.RunningTasksCount) > 0
		r := &models.Resource{
			ID:     aws.StringValue(instance.Ec2InstanceId),
			InUse:  inUse,
			Ports:  usedPorts,
			Memory: availableMemory.Megabytes(),
			CPU:    availableCPU.Megabytes(),
		}

		result = append(result, r)
	}

	return result, nil
}

func (e *EnvironmentScaler) getResourceConsumers(clusterName string) ([]*models.Resource, error) {
	// from scaler in develop, there are funcs to aggregate three types of resources
	//   - getPendingServiceResources
	//   - getPendingTaskResourcesInECS
	//   - getPendingTaskResourcesInJobs

	// GET PENDING SERVICE RESOURCE CONSUMERS
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
		// https://github.com/aws/aws-sdk-go/blob/ee1f179877b2daf2aaabf71fa900773bf8842253/service/ecs/api.go#L5420
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

	// TODO: can we use the service provider?

	resourceConsumers := []models.ResourceConsumer{}
	// TODO: consumers should have ports, mem, cpu

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

	// GET PENDING TASK RESOURCE CONSUMERS IN ECS
	taskARNs := []string{}
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

	// GET PENDING TASK RESOURCE CONSUMERS IN JOBS
	jobs, err := e.JobStore.SelectAll()
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		if job.Type == models.CreateTaskJob {
			// TODO: maybe remove pending check here
			if job.Status == models.PendingJobStatus || job.Status == models.InProgressJobStatus {
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

func (e *EnvironmentScaler) calculateNewProvider(clusterName string) (*models.Resource, error) {
	env, err := e.EnvironmentProvider.Read(clusterName)
	if err != nil {
		return nil, err
	}

	memory, err := strconv.Atoi(env.InstanceSize)
	if err != nil {
		return nil, err
	}

	// TODO: Calculate CPU
	resource := &models.Resource{}
	resource.ID = "<new instance>"
	resource.InUse = false
	resource.Ports = defaultPorts
	resource.Memory = float64(memory)

	return resource, nil
}

func (e *EnvironmentScaler) getContainerResourceFromDeploy(deployID string) ([]*models.Resource, error) {
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

	consumers := make([]*models.Resource, len(output.TaskDefinition.ContainerDefinitions))
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

		consumers[i] = models.ResourceConsumer{
			// TODO: add cpu to this model
			ID:     "",
			Memory: fmt.Sprintf("%v", memory), //todo: translate bytesize.ByteSize to string correctly
			// TODO: use bytesize object
			Ports: ports,
		}
	}

	// FIX: assignment to entry in nil map
	// e.deployCache[deployID] = consumers
	return consumers, nil
}

func (e *EnvironmentScaler) scaleTo(environmentID string, scale int, unusedProviders []*models.Resource) (int, error) {
	clusterName := addLayer0Prefix(e.Config.Instance(), environmentID)

	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	input.SetAutoScalingGroupNames([]*string{&clusterName})

	asg, err := e.Client.AutoScaling.DescribeAutoScalingGroups(input)
	if err != nil {
		return 0, err
	}

	currentCapacity := int(aws.Int64Value(asg.AutoScalingGroups[0].DesiredCapacity))

	switch {
	case scale > currentCapacity:
		log.Printf("Environment %s is attempting to scale up to size %d", clusterName, scale)
		return e.scaleUp(clusterName, scale, asg.AutoScalingGroups[0])
	case scale < currentCapacity:
		log.Printf("Environment %s is attempting to scale down to size %d", clusterName, scale)
		return e.scaleDown(clusterName, scale, asg.AutoScalingGroups[0], unusedProviders)
	default:
		log.Printf("Environment %s is at desired scale of %d. No scaling action required.", clusterName, scale)
		return currentCapacity, nil
	}
}

func (e *EnvironmentScaler) scaleUp(ecsEnvironmentID string, scale int, asg *autoscaling.Group) (int, error) {
	maxCapacity := int(aws.Int64Value(asg.MaxSize))
	if scale > maxCapacity {
		log.Printf("Scale %d is above the maximum capacity of %d. Setting desired capacity to %d.", scale, maxCapacity, maxCapacity)
		input := &autoscaling.UpdateAutoScalingGroupInput{}
		input.SetAutoScalingGroupName(aws.StringValue(asg.AutoScalingGroupName))
		input.SetMaxSize(int64(scale))

		if _, err := e.Client.AutoScaling.UpdateAutoScalingGroup(input); err != nil {
			return 0, err
		}
	}

	input := &autoscaling.SetDesiredCapacityInput{}
	input.SetAutoScalingGroupName(aws.StringValue(asg.AutoScalingGroupName))
	input.SetDesiredCapacity(int64(scale))

	if _, err := e.Client.AutoScaling.SetDesiredCapacity(input); err != nil {
		return 0, err
	}

	return scale, nil
}

func (e *EnvironmentScaler) scaleDown(ecsEnvironmentID string, scale int, asg *autoscaling.Group, unusedProviders []*models.Resource) (int, error) {
	minCapacity := int(aws.Int64Value(asg.MinSize))
	if scale < minCapacity {
		log.Printf("Scale %d is below the minimum capacity of %d. Setting desired capacity to %d.", scale, minCapacity, minCapacity)
		scale = minCapacity
	}

	currentCapacity := int(aws.Int64Value(asg.DesiredCapacity))
	if scale == currentCapacity {
		log.Printf("Environment %s is at desired scale of %d. No scaling action required.", ecsEnvironmentID, scale)
		return scale, nil
	}

	if scale < currentCapacity {
		input := &autoscaling.SetDesiredCapacityInput{}
		input.SetAutoScalingGroupName(aws.StringValue(asg.AutoScalingGroupName))
		input.SetDesiredCapacity(int64(scale))

		if _, err := e.Client.AutoScaling.SetDesiredCapacity(input); err != nil {
			return 0, err
		}
	}

	// choose which instances to terminate during our scale down process
	// instead of having asg randomly selecting instances
	// e.g. if we scale from 5->3, we can terminate up to 2 unused instances
	maxNumberOfInstancesToTerminate := currentCapacity - scale

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
		log.Printf("Environment %s terminating unused instance '%s'", ecsEnvironmentID, unusedProvider.ID)

		input := &autoscaling.TerminateInstanceInAutoScalingGroupInput{}
		input.SetInstanceId(unusedProvider.ID)
		input.SetShouldDecrementDesiredCapacity(false)

		if _, err := e.Client.AutoScaling.TerminateInstanceInAutoScalingGroup(input); err != nil {
			return 0, err
		}
	}

	return scale, nil
}

func (e *EnvironmentScaler) determineSortPriority(providers []*models.Resource, consumers []*models.Resource, clusterName string) int {
	cpuInstanceSize := len(providers)
	memoryInstanceSize := len(providers)

	for sortMethod := range []int{CPU, Memory} {

		for _, consumer := range consumers {
			hasRoom := false

			switch sortMethod {
			case CPU:
				sortByCPU(providers, consumers)
			case Memory:
				sortByMemory(providers, consumers)
			}

			// next, place any unused providers in the back of the list
			// that way, we can can delete them if we avoid placing any tasks in them
			sortByUsage(providers)

			for _, provider := range providers {
				if hasResourcesFor(*consumer, *provider) {
					hasRoom = true
					subtractResourcesFor(consumer, provider)
					break
				}
			}

			if !hasRoom {
				newProvider, err := e.calculateNewProvider(clusterName)
				// MARK Fix
				if err != nil {
					log.Println("ERROR")
					continue
				}
				if !hasResourcesFor(*consumer, *newProvider) {
					continue
				}

				switch sortMethod {
				case CPU:
					cpuInstanceSize++
				case Memory:
					memoryInstanceSize++
				}
			}
		}
	}

	if cpuInstanceSize > memoryInstanceSize {
		return CPU
	}

	return Memory
}

func hasResourcesFor(consumer models.Resource, provider models.Resource) bool {
	for _, wanted := range consumer.Ports {
		for _, used := range provider.Ports {
			if wanted == used {
				return false
			}
		}
	}

	return consumer.Memory <= provider.Memory
}

func subtractResourcesFor(consumer *models.Resource, provider *models.Resource) error {
	if !hasResourcesFor(*consumer, *provider) {
		return errors.Newf(errors.InvalidRequest, "Provider does not have adequate resources to subtract")
	}

	provider.Ports = append(provider.Ports, consumer.Ports...)
	provider.Memory -= consumer.Memory
	provider.CPU -= consumer.CPU
	provider.InUse = true

	return nil
}

func sortByMemory(resources ...[]*models.Resource) {
	for _, r := range resources {
		sorter := &ResourceSorter{
			Providers: r,
			lessThan: func(i *models.Resource, j *models.Resource) bool {
				return i.Memory < j.Memory
			},
		}

		sort.Sort(sorter)
	}
}

func sortByCPU(resources ...[]*models.Resource) {
	for _, r := range resources {
		sorter := &ResourceSorter{
			Providers: r,
			lessThan: func(i *models.Resource, j *models.Resource) bool {
				return i.CPU < j.CPU
			},
		}

		sort.Sort(sorter)
	}
}

func sortByUsage(r []*models.Resource) {
	sorter := &ResourceSorter{
		Providers: r,
		lessThan: func(i *models.Resource, j *models.Resource) bool {
			return i.InUse && !j.InUse
		},
	}

	sort.Sort(sorter)
}

type ResourceSorter struct {
	Providers []*models.Resource
	lessThan  func(*models.Resource, *models.Resource) bool
}

func (r *ResourceSorter) Len() int {
	return len(r.Providers)
}

func (r *ResourceSorter) Swap(i, j int) {
	r.Providers[i], r.Providers[j] = r.Providers[j], r.Providers[i]
}

func (r *ResourceSorter) Less(i, j int) bool {
	return r.lessThan(r.Providers[i], r.Providers[j])
}
