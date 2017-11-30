package aws

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/zpatrick/go-bytesize"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/api/provider"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/api/job"
)

type EnvironmentScaler struct {
	Client              *awsc.Client
	EnvironmentProvider provider.EnvironmentProvider
	ServiceProvider     provider.ServiceProvider
	TaskProvider        provider.TaskProvider
	JobStore    				job.Store
	Config              config.APIConfig
	deployCache         map[string][]models.ResourceConsumer
}

func NewEnvironmentScaler(a *awsc.Client, e provider.EnvironmentProvider, s provider.ServiceProvider, t provider.TaskProvider, j job.Store, c config.APIConfig) *EnvironmentScaler {
	return &EnvironmentScaler{
		Client:              a,
		EnvironmentProvider: e,
		ServiceProvider:     s,
		TaskProvider:        t,
		JobStore:     j,
		Config:              c,
	}
}

func (e *EnvironmentScaler) Scale(environmentID string) error {
	// GET RESOURCE PROVIDERS
	clusterName := addLayer0Prefix(e.Config.Instance(), environmentID)

	resourceProviders, err := e.getResourceProviders(clusterName)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] resourceProviders for env '%s': %#v", environmentID, resourceProviders)

	// GET RESOURCE CONSUMERS
	resourceConsumers, err := e.getResourceConsumers(clusterName)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] resourceConsumers for env '%s': %#v", environmentID, resourceConsumers)

	// e.scale(clusterName, resourceProviders, resourceConsumers)

	return fmt.Errorf("EnvironmentScaler not implemented")
}

func (e *EnvironmentScaler) scale(clusterName string, providers []*models.ResourceProvider, consumers []models.ResourceConsumer) (*models.ScalerRunInfo, error) {
	scaleBeforeRun := len(providers)
	var errs []error

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
			newProvider := &models.ResourceProvider{ID: clusterName}

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
	unusedProviders := []*models.ResourceProvider{}
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

	info := &models.ScalerRunInfo{
		EnvironmentID:           clusterName,
		PendingResources:        resourceConsumerModels(consumers),
		ResourceProviders:       resourceProviderModels(providers),
		ScaleBeforeRun:          scaleBeforeRun,
		DesiredScaleAfterRun:    desiredScale,
		ActualScaleAfterRun:     actualScale,
		UnusedResourceProviders: len(unusedProviders),
	}

	return info, errors.MultiError(errs)
}

func resourceConsumerModels(consumers []models.ResourceConsumer) []models.ResourceConsumer {
	consumerModels := make([]models.ResourceConsumer, len(consumers))
	for i, consumer := range consumers {
		newConsumer := models.ResourceConsumer{
			ID:     consumer.ID,
			Memory: consumer.Memory,
			Ports:  consumer.Ports,
		}
		consumerModels[i] = newConsumer
	}

	return consumerModels
}

func resourceProviderModels(providers []*models.ResourceProvider) []models.ResourceProvider {
	providerModels := make([]models.ResourceProvider, len(providers))
	for i, provider := range providers {
		newProvider := models.ResourceProvider{
			ID:              provider.ID,
			InUse:           provider.InUse,
			UsedPorts:       provider.UsedPorts,
			AvailableMemory: provider.AvailableMemory,
		}
		providerModels[i] = newProvider
	}

	return providerModels
}

func (e *EnvironmentScaler) getResourceProviders(clusterName string) ([]models.ResourceProvider, error) {
	listContainerInstancesInput := &ecs.ListContainerInstancesInput{}
	listContainerInstancesInput.SetCluster(clusterName)

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

	result := []models.ResourceProvider{}
	if len(output.ContainerInstances) == 0 {
		return result, nil
	}

	for _, instance := range output.ContainerInstances {
		if !aws.BoolValue(instance.AgentConnected) {
			continue
		}

		if aws.StringValue(instance.Status) != "ACTIVE" {
			continue
		}

		// it's non-intuitive, but the ports being used by the tasks live in
		// instance.RemainingResources, not instance.RegisteredResources
		var (
			usedPorts       []int
			availableMemory bytesize.Bytesize
		)

		for _, resource := range instance.RemainingResources {
			switch aws.StringValue(resource.Name) {
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
			}
		}

		inUse := aws.Int64Value(instance.PendingTasksCount)+aws.Int64Value(instance.RunningTasksCount) > 0
		r := models.ResourceProvider{
			ID:              aws.StringValue(instance.Ec2InstanceId),
			InUse:           inUse,
			UsedPorts:       usedPorts,
			AvailableMemory: fmt.Sprintf("%v", availableMemory), //todo: this is definitely not correct
		}

		result = append(result, r)
	}

	return result, nil
}

func (e *EnvironmentScaler) getResourceConsumers(clusterName string) ([]models.ResourceConsumer, error) {
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

	result := []models.ResourceConsumer{}

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

			result = append(result, c...)
		}
	}

	// GET PENDING TASK RESOURCE CONSUMERS IN ECS
	taskSummaries, err := e.TaskProvider.List()
	if err != nil {
		return nil, err
	}

	resourceConsumers := []models.ResourceConsumer{}

	for _, summary := range taskSummaries {
		if summary.EnvironmentID == clusterName {
			
			task, err := e.TaskProvider.Read(summary.TaskID)
			if err != nil {
				return nil, err
			}


			if task.PendingCount == 0 {
				continue
			}

			deployIDCopies := map[string]int{
				task.DeployID: int(task.TaskID)
			}

			// resource consumer ids are just used for debugging purposes
			generateID := func(deployID, containerName string, copy int) string {
				return fmt.Sprintf("Task: %s, Deploy: %s, Container: %s, Copy: %d", summary.TaskID, deployID, containerName, copy)
			}

			taskResourceConsumers, err := c.getResourcesHelper(deployIDCopies, generateID)
			if err != nil {
				return nil, err
			}

			resourceConsumers = append(resourceConsumers, taskResourceConsumers...)
		}
	}

	// return resourceConsumers, nil
	// input := &ecs.DescribeTasksInput{}
	// input.SetCluster(clusterName)
	// output, err := e.Client.ECS.DescribeTasks(input)
	// if err != nil {
	// 	return nil, err
	// }

	// for _, task := range output.Tasks {
	// }

	// GET PENDING TASK RESOURCE CONSUMERS IN JOBS



	jobs, err := e.JobStore.SelectAll()
	if err != nil {
		return nil, err
	}

	resourceConsumersJob := []models.ResourceConsumer{}
	for _, job := range jobs {
		if job.Type == int64(type.CreateTaskJob) {
			if job.JobStatus == int64(types.Pending) || job.JobStatus == int64(types.InProgress) {
				var req models.CreateTaskRequest
				if err := json.Unmarshal([]byte(job.Request), &req); err != nil {
					return nil, err
				}

				if req.EnvironmentID == environmentID {
					// note that this isn't exact if the job has started some, but not all of the tasks
					deployIDCopies := map[string]int{
						req.DeployID: int(req.Copies),
					}

					// resource consumer ids are just used for debugging purposes
					generateID := func(deployID, containerName string, copy int) string {
						return fmt.Sprintf("Task: %s, Deploy: %s, Container: %s, Copy: %d", req.TaskName, deployID, containerName, copy)
					}

					taskResourceConsumers, err := c.getResourcesHelper(deployIDCopies, generateID)
					if err != nil {
						return nil, err
					}

					resourceConsumers = append(resourceConsumers, taskResourceConsumers...)
				}
			}
		}
	}

	return resourceConsumers, nil


	return nil, nil
}

func (e *EnvironmentScaler) calculateNewProvider(clusterName string) (*models.ResourceProvider, error) {
	input := &autoscaling.CreateAutoScalingGroupInput{}
	if err := input.Validate(); err != nil {
		return nil, err
	}

	group, err := e.Client.AutoScaling.CreateAutoScalingGroup(input)
	if err != nil {
		return nil, err
	}

	inputLaunch := &autoscaling.DescribeLaunchConfigurationsInput{}
	inputLaunch.SetLaunchConfigurationNames([]*string{aws.String(group.String())})

	config, err := e.Client.AutoScaling.DescribeLaunchConfigurations(inputLaunch)
	if err != nil {
		return nil, err
	}

	_ = config
	// memory, ok := ec2.InstanceSizes[pstring(config.InstanceType)]
	// if !ok {
	// 	return nil, fmt.Errorf("Environment %s is using unknown instance type '%s'", environmentID, pstring(config.InstanceType))
	// }

	// these ports are automatically used by the ecs agent
	defaultPorts := []int{
		22,
		2376,
		2375,
		51678,
		51679,
	}

	resource := &models.ResourceProvider{}
	resource.ID = "<new instance>"
	resource.InUse = false
	resource.UsedPorts = defaultPorts
	// resource.AvailableMemory = memory

	return resource, nil
}

func (e *EnvironmentScaler) getContainerResourceFromDeploy(deployID string) ([]models.ResourceConsumer, error) {
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

	consumers := make([]models.ResourceConsumer, len(output.TaskDefinition.ContainerDefinitions))
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
			ID:     "",
			Memory: fmt.Sprintf("%v", memory), //todo: translate bytesize.ByteSize to string correctly
			Ports:  ports,
		}
	}

	e.deployCache[deployID] = consumers
	return consumers, nil
}

func (e *EnvironmentScaler) scaleTo(environmentID string, scale int, unusedProviders []*models.ResourceProvider) (int, error) {
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

func (e *EnvironmentScaler) scaleDown(ecsEnvironmentID string, scale int, asg *autoscaling.Group, unusedProviders []*models.ResourceProvider) (int, error) {
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

func sortProvidersByMemory(p []*models.ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: p,
		lessThan: func(i *models.ResourceProvider, j *models.ResourceProvider) bool {
			return i.AvailableMemory < j.AvailableMemory
		},
	}

	sort.Sort(sorter)
}

func sortProvidersByUsage(p []*models.ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: p,
		lessThan: func(i *models.ResourceProvider, j *models.ResourceProvider) bool {
			return i.InUse && !j.InUse
		},
	}

	sort.Sort(sorter)
}

func hasResourcesFor(consumer models.ResourceConsumer, provider *models.ResourceProvider) bool {
	for _, wanted := range consumer.Ports {
		for _, used := range provider.UsedPorts {
			if wanted == used {
				return false
			}
		}
	}

	return consumer.Memory <= provider.AvailableMemory
}

func subtractResourcesFor(consumer models.ResourceConsumer, provider *models.ResourceProvider) error {
	if !hasResourcesFor(consumer, provider) {
		return errors.Newf(errors.InvalidRequest, "Provider does not have adequate resources to subtract")
	}

	provider.UsedPorts = append(provider.UsedPorts, consumer.Ports...)

	// MARK: Strings
	// provider.AvailableMemory -= consumer.Memory
	provider.InUse = true

	return nil
}

type ResourceProviderSorter struct {
	Providers []*models.ResourceProvider
	lessThan  func(*models.ResourceProvider, *models.ResourceProvider) bool
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
