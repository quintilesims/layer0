package aws

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"

	"gitlab.imshealth.com/xfra/layer0/common/aws/ec2"

	"github.com/zpatrick/go-bytesize"

	"github.com/aws/aws-sdk-go/aws"
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
	deployCache         map[string][]models.ResourceConsumer
}

func NewEnvironmentScaler(a *awsc.Client, e provider.EnvironmentProvider, s provider.ServiceProvider, t provider.TaskProvider, j job.Store, c config.APIConfig) *EnvironmentScaler {
	return &EnvironmentScaler{
		Client:              a,
		EnvironmentProvider: e,
		ServiceProvider:     s,
		TaskProvider:        t,
		JobStore:            j,
		Config:              c,
	}
}

func (e *EnvironmentScaler) Scale(environmentID string) error {
	clusterName := addLayer0Prefix(e.Config.Instance(), environmentID)

	resourceProviders, err := e.getResourceProviders(clusterName)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] resourceProviders for env '%s': %#v", environmentID, resourceProviders)

	var resourceConsumers []models.ResourceConsumer

	// get pending service resource consumers
	serviceConsumers, err := e.getResourceConsumers_PendingServices(clusterName)
	if err != nil {
		return err
	}

	resourceConsumers = append(resourceConsumers, serviceConsumers...)

	// get pending task resource consumers in ECS
	ecsTaskConsumers, err := e.getResourceConsumers_TasksInECS(clusterName)
	if err != nil {
		return err
	}

	resourceConsumers = append(resourceConsumers, ecsTaskConsumers...)

	// get pending task resource consumers in jobs
	jobTaskConsumers, err := e.getResourceConsumers_TasksInJobs(clusterName)
	if err != nil {
		return err
	}

	resourceConsumers = append(resourceConsumers, jobTaskConsumers...)

	log.Printf("[DEBUG] resourceConsumers for env '%s': %#v", environmentID, resourceConsumers)

	var errs []error

	// calculate/check for scaling up
	resourceProviders, scalingErrs := e.calculateScaleUp(clusterName, resourceProviders, resourceConsumers)

	errs = append(errs, scalingErrs...)

	// calculate/check for scaling down

	// do the scaling

	return nil
}

func (e *EnvironmentScaler) calculateNewProvider(clusterName string) (*models.ResourceProvider, error) {
	env, err := e.EnvironmentProvider.Read(clusterName)
	if err != nil {
		return nil, err
	}

	resource := &models.ResourceProvider{}
	resource.AvailableMemory = ec2.InstanceSizes[env.InstanceSize].(bytesize.Bytesize)
	resource.ID = "<new instance>"
	resource.InUse = false
	resource.UsedPorts = defaultPorts

	return resource, nil
}

func (e *EnvironmentScaler) calculateScaleUp(clusterName string, resourceProviders []*models.ResourceProvider, resourceConsumers []models.ResourceConsumer) ([]*models.ResourceProvider, []error) {
	sorts := []struct {
		Errs      []error
		Providers []*models.ResourceProvider
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
				if hasResourcesFor(provider, consumer) {
					hasRoom = true
					subtractResourcesFor(provider, consumer)
					break
				}
			}

			if !hasRoom {
				newProvider, err := e.calculateNewProvider(clusterName)
				if err != nil {
					s.Errs = append(s.Errs, err)
					continue
				}

				if !hasResourcesFor(newProvider, consumer) {
					text := fmt.Sprintf("Resource '%s' cannot fit into an empty provider!", consumer.ID)
					text += "\nThe instance size in your environment is too small to run this resource."
					text += "\n Pleace increase the instance size for your environment."
					err := fmt.Errorf(text)
					s.Errs = append(s.Errs, err)
					continue
				}

				subtractResourcesFor(newProvider, consumer)
				s.Providers = append(s.Providers, newProvider)
			}
		}
	}

	if len(sorts[0].Providers) > len(sorts[1].Providers) {
		return sorts[0].Providers, sorts[0].Errs
	}

	return sorts[1].Providers, sorts[1].Errs
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
			Memory: memory,
			Ports:  ports,
		}
	}

	e.deployCache[deployID] = consumers
	return consumers, nil
}

func (e *EnvironmentScaler) getResourceConsumers_PendingServices(clusterName string) ([]models.ResourceConsumer, error) {
	var resourceConsumers []models.ResourceConsumer

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

func (e *EnvironmentScaler) getResourceConsumers_TasksInECS(clusterName string) ([]models.ResourceConsumer, error) {
	var (
		taskARNs          []string
		resourceConsumers []models.ResourceConsumer
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

func (e *EnvironmentScaler) getResourceConsumers_TasksInJobs(clusterName string) ([]models.ResourceConsumer, error) {
	var resourceConsumers []models.ResourceConsumer
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

func (e *EnvironmentScaler) getResourceProviders(clusterName string) ([]*models.ResourceProvider, error) {
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

	result := []*models.ResourceProvider{}
	if len(containerInstanceARNs) == 0 {
		return result, nil
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
						log.Printf("[WARN] Instance %s: Failed to convert port to int: %v", aws.StringValue(instance.Ec2InstanceId), err)
						continue
					}

					usedPorts = append(usedPorts, port)
				}
			}
		}

		inUse := aws.Int64Value(instance.PendingTasksCount)+aws.Int64Value(instance.RunningTasksCount) > 0
		r := &models.ResourceProvider{
			AgentConnected:  aws.BoolValue(instance.AgentConnected),
			AvailableCPU:    availableCPU,
			AvailableMemory: availableMemory,
			ID:              aws.StringValue(instance.Ec2InstanceId),
			InUse:           inUse,
			Status:          aws.StringValue(instance.Status),
			UsedPorts:       usedPorts,
		}

		result = append(result, r)
	}

	return result, nil
}

func hasResourcesFor(provider *models.ResourceProvider, consumer models.ResourceConsumer) bool {
	for _, wanted := range consumer.Ports {
		for _, used := range provider.UsedPorts {
			if wanted == used {
				return false
			}
		}
	}

	return consumer.CPU <= provider.AvailableCPU && consumer.Memory <= provider.AvailableMemory
}

func sortConsumersByCPU(c []models.ResourceConsumer) {
	sorter := &ResourceConsumerSorter{
		Consumers: c,
		lessThan: func(i models.ResourceConsumer, j models.ResourceConsumer) bool {
			return i.CPU < j.CPU
		},
	}

	sort.Sort(sorter)
}

func sortConsumersByMemory(c []models.ResourceConsumer) {
	sorter := &ResourceConsumerSorter{
		Consumers: c,
		lessThan: func(i models.ResourceConsumer, j models.ResourceConsumer) bool {
			return i.Memory < j.Memory
		},
	}

	sort.Sort(sorter)
}

func sortProvidersByCPU(p []*models.ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: p,
		lessThan: func(i *models.ResourceProvider, j *models.ResourceProvider) bool {
			return i.AvailableCPU < j.AvailableCPU
		},
	}

	sort.Sort(sorter)
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

func subtractResourcesFor(provider *models.ResourceProvider, consumer models.ResourceConsumer) error {
	if !hasResourcesFor(provider, consumer) {
		return errors.Newf(errors.InvalidRequest, "Provider does not have adequate resources to subtract.")
	}

	provider.AvailableCPU -= consumer.CPU
	provider.AvailableMemory -= consumer.Memory
	provider.InUse = true
	provider.UsedPorts = append(provider.UsedPorts, consumer.Ports...)

	return nil
}

type ResourceConsumerSorter struct {
	Consumers []models.ResourceConsumer
	lessThan  func(models.ResourceConsumer, models.ResourceConsumer) bool
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
