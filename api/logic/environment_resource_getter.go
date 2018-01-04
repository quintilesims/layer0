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

type EnvironmentResourceGetter struct {
	ServiceLogic ServiceLogic
	TaskLogic    TaskLogic
	DeployLogic  DeployLogic
	JobLogic     JobLogic
	deployCache  map[string][]resource.ResourceConsumer
}

func NewEnvironmentResourceGetter(s ServiceLogic, t TaskLogic, d DeployLogic, j JobLogic) *EnvironmentResourceGetter {
	return &EnvironmentResourceGetter{
		ServiceLogic: s,
		TaskLogic:    t,
		DeployLogic:  d,
		JobLogic:     j,
		deployCache:  map[string][]resource.ResourceConsumer{},
	}
}

func (c *EnvironmentResourceGetter) GetConsumers(environmentID string) ([]resource.ResourceConsumer, error) {
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

func (c *EnvironmentResourceGetter) getPendingServiceResources(environmentID string) ([]resource.ResourceConsumer, error) {
	resourceConsumers := []resource.ResourceConsumer{}
	services, err := c.ServiceLogic.GetEnvironmentServices(environmentID)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		deployIDCopies := map[string]int{}
		for _, deployment := range service.Deployments {
			// deployment.RunningCount is the number of containers already running on an instance
			// deployment.PendingCount is the number of containers that are alraedy on an instance, but are being pulled
			// we only care about containers that are not on instances yet

			if numPending := deployment.DesiredCount - (deployment.RunningCount + deployment.PendingCount); numPending > 0 {
				deployIDCopies[deployment.DeployID] = int(numPending)
			}
		}

		if len(deployIDCopies) == 0 {
			continue
		}

		// resource consumer ids are just used for debugging purposes
		generateID := func(deployID, containerName string, copy int) string {
			return fmt.Sprintf("Service: %s, Deploy: %s, Container: %s, Copy: %d", service.ServiceID, deployID, containerName, copy)
		}

		serviceResourceConsumers, err := c.getResourcesHelper(deployIDCopies, generateID)
		if err != nil {
			return nil, err
		}

		resourceConsumers = append(resourceConsumers, serviceResourceConsumers...)
	}

	return resourceConsumers, nil
}

func (c *EnvironmentResourceGetter) getPendingTaskResourcesInECS(environmentID string) ([]resource.ResourceConsumer, error) {
	tasks, err := c.TaskLogic.GetEnvironmentTasks(environmentID)
	if err != nil {
		return nil, err
	}

	resourceConsumers := []resource.ResourceConsumer{}
	for _, task := range tasks {
		if task.PendingCount == 0 {
			continue
		}

		deployIDCopies := map[string]int{
			task.DeployID: int(task.PendingCount),
		}

		// resource consumer ids are just used for debugging purposes
		generateID := func(deployID, containerName string, copy int) string {
			return fmt.Sprintf("Task: %s, Deploy: %s, Container: %s, Copy: %d", task.TaskID, deployID, containerName, copy)
		}

		taskResourceConsumers, err := c.getResourcesHelper(deployIDCopies, generateID)
		if err != nil {
			return nil, err
		}

		resourceConsumers = append(resourceConsumers, taskResourceConsumers...)
	}

	return resourceConsumers, nil
}

func (c *EnvironmentResourceGetter) getPendingTaskResourcesInJobs(environmentID string) ([]resource.ResourceConsumer, error) {
	jobs, err := c.JobLogic.ListJobs()
	if err != nil {
		return nil, err
	}

	resourceConsumers := []resource.ResourceConsumer{}
	for _, job := range jobs {
		if job.JobType == int64(types.CreateTaskJob) {
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
}

func (c *EnvironmentResourceGetter) getResourcesHelper(deployIDCopies map[string]int, generateID func(string, string, int) string) ([]resource.ResourceConsumer, error) {
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

func (c *EnvironmentResourceGetter) getContainerResourcesFromDeploy(deployID string) ([]resource.ResourceConsumer, error) {
	if consumers, ok := c.deployCache[deployID]; ok {
		return consumers, nil
	}

	d, err := c.DeployLogic.GetDeploy(deployID)
	if err != nil {
		return nil, err
	}

	deploy, err := ecsbackend.MarshalDockerrun(d.Dockerrun)
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
