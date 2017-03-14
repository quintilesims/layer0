package logic

import (
	"encoding/json"
	"fmt"
	"github.com/quintilesims/layer0/api/backend/ecs"
	"github.com/quintilesims/layer0/api/scheduler/resource"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
	"github.com/zpatrick/go-bytesize"
)

// todo: caching
type ClusterResourceGetter struct {
	ServiceLogic ServiceLogic
	TaskLogic    TaskLogic
	DeployLogic  DeployLogic
	JobLogic     JobLogic
	deployCache  map[string][]resource.ResourceConsumer
}

func NewClusterResourceGetter(s ServiceLogic, t TaskLogic, d DeployLogic, j JobLogic) *ClusterResourceGetter {
	return &ClusterResourceGetter{
		ServiceLogic: s,
		TaskLogic:    t,
		DeployLogic:  d,
		JobLogic:     j,
		deployCache:  map[string][]resource.ResourceConsumer{},
	}
}

func (c *ClusterResourceGetter) GetPendingResources(environmentID string) ([]resource.ResourceConsumer, error) {
	serviceResources, err := c.getPendingServiceResources(environmentID)
	if err != nil {
		return nil, err
	}

	taskResourcesInECS, err := c.getPendingTaskResourcesInECS(environmentID)
	if err != nil {
		return nil, err
	}

	taskResourcesInJobs, err := c.getPendingTaskResourcesInJobs(environmentID)
	if err != nil {
		return nil, err
	}

	totalResources := append(serviceResources, taskResourcesInECS...)
	totalResources = append(totalResources, taskResourcesInJobs...)
	return totalResources, nil
}

func (c *ClusterResourceGetter) getPendingServiceResources(environmentID string) ([]resource.ResourceConsumer, error) {
	serviceSummaries, err := c.ServiceLogic.ListServices()
	if err != nil {
		return nil, err
	}

	resourceConsumers := []resource.ResourceConsumer{}
	for _, summary := range serviceSummaries {
		if summary.EnvironmentID == environmentID {
			service, err := c.ServiceLogic.GetService(summary.ServiceID)
			if err != nil {
				return nil, err
			}

			deployIDCopies := map[string]int{}
			for _, deployment := range service.Deployments {
				if numPending := int(deployment.DesiredCount - deployment.RunningCount); numPending > 0 {
					deployIDCopies[deployment.DeployID] = numPending
				}
			}

			// resource consumer ids are just used for debugging purposes
			generateID := func(deployID, containerName string, copy int) string {
				return fmt.Sprintf("Service: %s, Deploy: %s, Container: %s, Copy: %d", summary.ServiceID, deployID, containerName, copy)
			}

			serviceResourceConsumers, err := c.getResourcesHelper(deployIDCopies, generateID)
			if err != nil {
				return nil, err
			}

			resourceConsumers = append(resourceConsumers, serviceResourceConsumers...)
		}
	}

	return resourceConsumers, nil
}

func (c *ClusterResourceGetter) getPendingTaskResourcesInECS(environmentID string) ([]resource.ResourceConsumer, error) {
	taskSummaries, err := c.TaskLogic.ListTasks()
	if err != nil {
		return nil, err
	}

	resourceConsumers := []resource.ResourceConsumer{}
	for _, summary := range taskSummaries {
		if summary.EnvironmentID == environmentID {
			task, err := c.TaskLogic.GetTask(summary.TaskID)
			if err != nil {
				return nil, err
			}

			deployIDCopies := map[string]int{}
			if numPending := int(task.PendingCount); numPending > 0 {
				deployIDCopies[task.DeployID] = numPending
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

	return resourceConsumers, nil
}

func (c *ClusterResourceGetter) getPendingTaskResourcesInJobs(environmentID string) ([]resource.ResourceConsumer, error) {
	jobs, err := c.JobLogic.ListJobs()
	if err != nil {
		return nil, err
	}

	resourceConsumers := []resource.ResourceConsumer{}
	for _, job := range jobs {
		if job.JobType == int64(types.CreateTaskJob) {
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

	return resourceConsumers, nil
}

func (c *ClusterResourceGetter) getResourcesHelper(deployIDCopies map[string]int, generateID func(string, string, int) string) ([]resource.ResourceConsumer, error) {
	resourceConsumers := []resource.ResourceConsumer{}
	for deployID, copies := range deployIDCopies {
		containerResources, err := c.getContainerResourcesFromDeploy(deployID)
		if err != nil {
			return nil, err
		}

		for i := 0; i < copies; i++ {
			for _, containerResource := range containerResources {
				id := generateID(deployID, containerResource.ID, i+1)
				consumer := resource.NewResourceConsumer(id, containerResource.Memory, containerResource.Ports)
				resourceConsumers = append(resourceConsumers, consumer)
			}
		}
	}

	return resourceConsumers, nil
}

func (c *ClusterResourceGetter) getContainerResourcesFromDeploy(deployID string) ([]resource.ResourceConsumer, error) {
	if consumers, ok := c.deployCache[deployID]; ok {
		return consumers, nil
	}

	d, err := c.DeployLogic.GetDeploy(deployID)
	if err != nil {
		return nil, err
	}

	deploy, err := ecsbackend.MarshalDeploy(d.Dockerrun)
	if err != nil {
		return nil, err
	}

	consumers := make([]resource.ResourceConsumer, len(deploy.ContainerDefinitions))
	for i, container := range deploy.ContainerDefinitions {
		var memory bytesize.Bytesize

		if container.MemoryReservation != nil && *container.MemoryReservation != 0 {
			memory = bytesize.MiB * bytesize.Bytesize(*container.MemoryReservation)
		}

		if container.Memory != nil && *container.Memory != 0 {
			memory = bytesize.MiB * bytesize.Bytesize(*container.Memory)
		}

		ports := []int{}
		for _, p := range container.PortMappings {
			if p.HostPort != nil && *p.HostPort != 0 {
				ports = append(ports, int(*p.HostPort))
			}
		}

		consumers[i] = resource.NewResourceConsumer(*container.Name, memory, ports)
	}

	c.deployCache[deployID] = consumers
	return consumers, nil
}
