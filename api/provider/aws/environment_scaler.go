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

// https://aws.amazon.com/ec2/instance-types/
var InstanceSizes = map[string]bytesize.Bytesize{
	"t2.nano":     0.5 * bytesize.GiB,
	"t2.micro":    1 * bytesize.GiB,
	"t2.small":    2 * bytesize.GiB,
	"t2.medium":   4 * bytesize.GiB,
	"t2.large":    8 * bytesize.GiB,
	"m4.large":    8 * bytesize.GiB,
	"m4.xlarge":   16 * bytesize.GiB,
	"m4.2xlarge":  32 * bytesize.GiB,
	"m4.4xlarge":  64 * bytesize.GiB,
	"m4.10xlarge": 160 * bytesize.GiB,
	"m3.medium":   3.75 * bytesize.GiB,
	"m3.large":    7.5 * bytesize.GiB,
	"m3.xlarge":   15 * bytesize.GiB,
	"m3.2xlarge":  30 * bytesize.GiB,
	"c4.large":    3.75 * bytesize.GiB,
	"c4.xlarge":   7.5 * bytesize.GiB,
	"c4.2xlarge":  15 * bytesize.GiB,
	"c4.4xlarge":  30 * bytesize.GiB,
	"c4.8xlarge":  60 * bytesize.GiB,
	"c3.large":    3.75 * bytesize.GiB,
	"c3.xlarge":   7.5 * bytesize.GiB,
	"c3.2xlarge":  15 * bytesize.GiB,
	"c3.4xlarge":  30 * bytesize.GiB,
	"c3.8xlarge":  60 * bytesize.GiB,
	"g2.2xlarge":  15 * bytesize.GiB,
	"g2.8xlarge":  60 * bytesize.GiB,
	"x1.32xlarge": 1952 * bytesize.GiB,
	"r3.large":    15.25 * bytesize.GiB,
	"r3.xlarge":   30.5 * bytesize.GiB,
	"r3.2xlarge":  61 * bytesize.GiB,
	"r3.4xlarge":  122 * bytesize.GiB,
	"r3.8xlarge":  244 * bytesize.GiB,
	"i3.large":    15.25 * bytesize.GiB,
	"i3.xlarge":   30.5 * bytesize.GiB,
	"i3.2xlarge":  61 * bytesize.GiB,
	"i3.4xlarge":  122 * bytesize.GiB,
	"i3.8xlarge":  244 * bytesize.GiB,
	"d2.xlarge":   30.5 * bytesize.GiB,
	"d2.2xlarge":  61 * bytesize.GiB,
	"d2.4xlarge":  122 * bytesize.GiB,
	"d2.8xlarge":  244 * bytesize.GiB,
}

type EnvironmentScaler struct {
	Client              *awsc.Client
	EnvironmentProvider provider.EnvironmentProvider
	ServiceProvider     provider.ServiceProvider
	TaskProvider        provider.TaskProvider
	JobStore            job.Store
	Config              config.APIConfig
	deployCache         map[string][]ResourceConsumer
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
		// deployCache: map[string][]*Resource,
	}
}

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

func (e *EnvironmentScaler) GetCurrentState(clusterName string) ([]ResourceConsumer, []*ResourceProvider, error) {
	resourceProviders, err := e.getResourceProviders(clusterName)
	if err != nil {
		return nil, nil, err
	}

	humanReadableProviders := ""
	for _, r := range resourceProviders {
		humanReadableProviders += fmt.Sprintf("\n    %#v", r)
	}

	log.Printf("[DEBUG] [EnvironmentScaler] resourceProviders for env '%s': %s", clusterName, humanReadableProviders)

	var resourceConsumers []ResourceConsumer

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

func (e *EnvironmentScaler) CalculateOptimizedState(clusterName string, resourceConsumers []ResourceConsumer, resourceProviders []*ResourceProvider) ([]*ResourceProvider, []*ResourceProvider, []error, error) {
	// calculate for scaling up
	resourceProviders, scaleUpErrs, err := e.calculateScaleUp(clusterName, resourceProviders, resourceConsumers)
	if err != nil {
		return nil, nil, nil, err
	}

	// calculate for scaling down
	unusedProviders := e.calculateScaleDown(clusterName, resourceProviders)

	return resourceProviders, unusedProviders, scaleUpErrs, nil
}

func (e *EnvironmentScaler) ScaleToState(clusterName string, desiredScale int, unusedProviders []*ResourceProvider) error {
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

func (e *EnvironmentScaler) calculateNewProvider(clusterName string) (*ResourceProvider, error) {
	env, err := e.EnvironmentProvider.Read(clusterName)
	if err != nil {
		return nil, err
	}

	resource := &ResourceProvider{}
	resource.AvailableMemory = InstanceSizes[env.InstanceSize]
	resource.ID = "<new instance>"
	resource.InUse = false
	resource.UsedPorts = defaultPorts

	return resource, nil
}

func (e *EnvironmentScaler) calculateScaleDown(clusterName string, resourceProviders []*ResourceProvider) []*ResourceProvider {
	var unusedProviders []*ResourceProvider

	for _, provider := range resourceProviders {
		if !provider.InUse {
			unusedProviders = append(unusedProviders, provider)
		}
	}

	return unusedProviders
}

func (e *EnvironmentScaler) calculateScaleUp(clusterName string, resourceProviders []*ResourceProvider, resourceConsumers []ResourceConsumer) ([]*ResourceProvider, []error, error) {
	sorts := []struct {
		Errs      []error
		Providers []*ResourceProvider
		SortBy    string
	}{
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

	if len(sorts[0].Providers) > len(sorts[1].Providers) {
		return sorts[0].Providers, sorts[0].Errs, nil
	}

	return sorts[1].Providers, sorts[1].Errs, nil
}

func (e *EnvironmentScaler) getContainerResourceFromDeploy(deployID string) ([]ResourceConsumer, error) {
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

	consumers := make([]ResourceConsumer, len(output.TaskDefinition.ContainerDefinitions))
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

		consumers[i] = ResourceConsumer{
			ID:     "",
			Memory: memory,
			Ports:  ports,
		}
	}

	e.deployCache[deployID] = consumers
	return consumers, nil
}

func (e *EnvironmentScaler) getResourceConsumers_PendingServices(clusterName string) ([]ResourceConsumer, error) {
	var resourceConsumers []ResourceConsumer

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

func (e *EnvironmentScaler) getResourceConsumers_TasksInECS(clusterName string) ([]ResourceConsumer, error) {
	var (
		taskARNs          []string
		resourceConsumers []ResourceConsumer
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

func (e *EnvironmentScaler) getResourceConsumers_TasksInJobs(clusterName string) ([]ResourceConsumer, error) {
	var resourceConsumers []ResourceConsumer
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

func (e *EnvironmentScaler) getResourceProviders(clusterName string) ([]*ResourceProvider, error) {
	var result []*ResourceProvider

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
			availableCPU    bytesize.Bytesize
			availableMemory bytesize.Bytesize
		)

		for _, resource := range instance.RemainingResources {
			switch aws.StringValue(resource.Name) {
			case "CPU":
				val := aws.Int64Value(resource.IntegerValue)
				availableCPU = bytesize.MiB * bytesize.Bytesize(val)

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

		r := &ResourceProvider{
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

func (e *EnvironmentScaler) scaleDown(clusterName string, desiredScale int, asg *autoscaling.Group, unusedProviders []*ResourceProvider) error {
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

func sortConsumersByCPU(c []ResourceConsumer) {
	sorter := &ResourceConsumerSorter{
		Consumers: c,
		lessThan: func(i ResourceConsumer, j ResourceConsumer) bool {
			return i.CPU < j.CPU
		},
	}

	sort.Sort(sorter)
}

func sortConsumersByMemory(c []ResourceConsumer) {
	sorter := &ResourceConsumerSorter{
		Consumers: c,
		lessThan: func(i ResourceConsumer, j ResourceConsumer) bool {
			return i.Memory < j.Memory
		},
	}

	sort.Sort(sorter)
}

func sortProvidersByCPU(p []*ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: p,
		lessThan: func(i *ResourceProvider, j *ResourceProvider) bool {
			return i.AvailableCPU < j.AvailableCPU
		},
	}

	sort.Sort(sorter)
}

func sortProvidersByMemory(p []*ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: p,
		lessThan: func(i *ResourceProvider, j *ResourceProvider) bool {
			return i.AvailableMemory < j.AvailableMemory
		},
	}

	sort.Sort(sorter)
}

func sortProvidersByUsage(r []*ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: r,
		lessThan: func(i *ResourceProvider, j *ResourceProvider) bool {
			return i.InUse && !j.InUse
		},
	}

	sort.Sort(sorter)
}

type ResourceConsumerSorter struct {
	Consumers []ResourceConsumer
	lessThan  func(ResourceConsumer, ResourceConsumer) bool
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
	Providers []*ResourceProvider
	lessThan  func(*ResourceProvider, *ResourceProvider) bool
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

type ResourceProvider struct {
	AgentIsConnected bool              `json:"agent_connected"`
	AvailableCPU     bytesize.Bytesize `json:"available_cpu"`
	AvailableMemory  bytesize.Bytesize `json:"available_memory"`
	ID               string            `json:"id"`
	InUse            bool              `json:"in_use"`
	Status           string            `json:"status"`
	UsedPorts        []int             `json:"used_ports"`
}

func (r *ResourceProvider) HasResourcesFor(consumer ResourceConsumer) bool {
	if !r.AgentIsConnected || r.Status != ecs.ContainerInstanceStatusActive {
		return false
	}

	for _, wanted := range consumer.Ports {
		for _, used := range r.UsedPorts {
			if wanted == used {
				return false
			}
		}
	}

	return consumer.CPU <= r.AvailableCPU && consumer.Memory <= r.AvailableMemory
}

func (r *ResourceProvider) SubtractResourcesFor(consumer ResourceConsumer) error {
	if !r.HasResourcesFor(consumer) {
		return fmt.Errorf("Cannot subtract resources for consumer '%s' from provider '%s'.", consumer.ID, r.ID)
	}

	r.AvailableCPU -= consumer.CPU
	r.AvailableMemory -= consumer.Memory
	r.InUse = true
	r.UsedPorts = append(r.UsedPorts, consumer.Ports...)

	return nil
}

type ResourceConsumer struct {
	CPU    bytesize.Bytesize `json:"cpu"`
	ID     string            `json:"id"`
	Memory bytesize.Bytesize `json:"memory"`
	Ports  []int             `json:"ports"`
}
