package aws

import (
	"fmt"
<<<<<<< HEAD
	"sort"
=======
	"log"
>>>>>>> @{-1}
	"strconv"

	"github.com/zpatrick/go-bytesize"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/api/provider"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type EnvironmentScaler struct {
	Client              *awsc.Client
	EnvironmentProvider provider.EnvironmentProvider
	ServiceProvider     provider.ServiceProvider
	Config              config.APIConfig
	deployCache         map[string][]models.ResourceConsumer
}

func NewEnvironmentScaler(a *awsc.Client, e provider.EnvironmentProvider, s provider.ServiceProvider, c config.APIConfig) *EnvironmentScaler {
	return &EnvironmentScaler{
		Client:              a,
		EnvironmentProvider: e,
		ServiceProvider:     s,
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

<<<<<<< HEAD
=======
	log.Printf("[DEBUG] resourceProviders: %#v", resourceProviders)

	// GET RESOURCE CONSUMERS
>>>>>>> @{-1}
	resourceConsumers, err := e.getResourceConsumers(clusterName)
	if err != nil {
		return err
	}
<<<<<<< HEAD
=======

	log.Printf("[DEBUG] resourceConsumers: %#v", resourceConsumers)

	// RUN BASIC SCALER
>>>>>>> @{-1}

	// calculate desired capacity

	// scale to new capacity if required

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

			newProvider.SubtractResourcesFor(consumer)
			providers = append(providers, newProvider)
		}
	}

	// check if we need to scale down
	unusedProviders := []*resource.ResourceProvider{}
	for i := 0; i < len(providers); i++ {
		if !providers[i].IsInUse() {
			unusedProviders = append(unusedProviders, providers[i])
		}
	}

	desiredScale := len(providers) - len(unusedProviders)
	actualScale, err := providerManager.ScaleTo(environmentID, desiredScale, unusedProviders)
	if err != nil {
		errs = append(errs, err)
	}

	info := &models.ScalerRunInfo{
		EnvironmentID:           environmentID,
		PendingResources:        resourceConsumerModels(consumers),
		ResourceProviders:       resourceProviderModels(providers),
		ScaleBeforeRun:          scaleBeforeRun,
		DesiredScaleAfterRun:    desiredScale,
		ActualScaleAfterRun:     actualScale,
		UnusedResourceProviders: len(unusedProviders),
	}

	return info, errors.MultiError(errs)
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

func (e *EnvironmentScaler) calculateNewProvider(clusterName string) (*models.ResourceProvider, error) {
	e.Client.CreateAutoScalingGroupInput
	group, e := e.Client.AutoScaling.CreateAutoScalingGroup()

	group, err := e.Client.AutoScaling.DescribeAutoScalingGroup(clusterName)
	if err != nil {
		return nil, err
	}

	config, err := r.Autoscaling.DescribeLaunchConfiguration(pstring(group.LaunchConfigurationName))
	if err != nil {
		return nil, err
	}

	memory, ok := ec2.InstanceSizes[pstring(config.InstanceType)]
	if !ok {
		return nil, fmt.Errorf("Environment %s is using unknown instance type '%s'", environmentID, pstring(config.InstanceType))
	}

	// these ports are automatically used by the ecs agent
	defaultPorts := []int{
		22,
		2376,
		2375,
		51678,
		51679,
	}

	resource := models.ResourceProvider{}
	resource.ID = "<new instance>"
	resource.InUse = false
	resource.UsedPorts = defaultPorts
	resource.AvailableMemory = memory

	return resource, nil
}

func (e *EnvironmentScaler) getResourceConsumers(clusterName string) ([]models.ResourceConsumer, error) {
	// from scaler in develop, there are funcs to aggregate three types of resources
	//   - getPendingServiceResources
	//   - getPendingTaskResourcesInECS
	//   - getPendingTaskResourcesInJobs

	// GET PENDING SERVICE RESOURCE CONSUMERS
	input := &ecs.DescribeServicesInput{}
	input.SetCluster(clusterName)
	output, err := e.Client.ECS.DescribeServices(input)
	if err != nil {
		return nil, err
	}

	result := []models.ResourceConsumer{}

	for _, service := range output.Services {
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

	// input := &ecs.DescribeTasksInput{}
	// input.SetCluster(clusterName)
	// output, err := e.Client.ECS.DescribeTasks(input)
	// if err != nil {
	// 	return nil, err
	// }

	// for _, task := range output.Tasks {
	// }

	// GET PENDING TASK RESOURCE CONSUMERS IN JOBS

	return nil, nil
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

// Helpers
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

func hasResourcesFor(consumer models.ResourceConsumer, provider *models.ResourceProvider) bool {
	for _, wanted := range consumer.Ports {
		for _, used := range provider.usedPorts {
			if wanted == used {
				return false
			}
		}
	}

	return consumer.Memory <= r.availableMemory
}

func subtractResourcesFor(consumer models.ResourceConsumer, provider *models.ResourceProvider) error {
	if !hasResourcesFor(consumer, provider) {
		return errors.New("Provider does not have adequate resources to subtract")
	}

	provider.UsedPorts = append(provider.UsedPorts, consumer.Ports...)
	provider.AvailableMemory -= consumer.Memory
	provider.InUse = true

	return nil
}
