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
}

func NewClusterResourceGetter(s ServiceLogic, t TaskLogic, d DeployLogic, j JobLogic) *ClusterResourceGetter {
	return &ClusterResourceGetter{
		ServiceLogic: s,
		TaskLogic:    t,
		DeployLogic:  d,
		JobLogic:     j,
	}
}

func (c *ClusterResourceGetter) GetPendingResources(environmentID string) ([]resource.ResourceConsumer, error) {
	serviceResources, err := c.getPendingServiceResources(environmentID)
	if err != nil {
		return nil, err
	}

	taskResources, err := c.getPendingTaskResources(environmentID)
	if err != nil {
		return nil, err
	}

	totalResources := append(serviceResources, taskResources...)
	return totalResources, nil
}

// todo: helper function for these, pass in:
// 	id generator
// 	map[Deploy]int = Deploy and how many copies
func (c *ClusterResourceGetter) getPendingServiceResources(environmentID string) ([]resource.ResourceConsumer, error) {
	services, err := c.ServiceLogic.ListServices()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(services); i++ {
		if services[i].EnvironmentID != environmentID {
			services = append(services[:i], services[i+1:]...)
			i--
		}
	}

	resourceConsumers := []resource.ResourceConsumer{}
	for _, service := range services {
		details, err := c.ServiceLogic.GetService(service.ServiceID)
		if err != nil {
			return nil, err
		}

		for _, deployment := range details.Deployments {
			if deployment.PendingCount > 0 {
				resources, err := c.getResourcesForDeploy(deployment.DeployID)
				if err != nil {
					return nil, err
				}

				for i := 0; i < int(deployment.PendingCount); i++ {
					for _, r := range resources {
						id := fmt.Sprintf("%s.%s.%s.%d", service.ServiceID, deployment.DeployID, r.ID, i)
						consumer := resource.NewResourceConsumer(id, r.Memory, r.Ports)
						resourceConsumers = append(resourceConsumers, consumer)
					}
				}
			}
		}
	}

	return resourceConsumers, nil
}

func (c *ClusterResourceGetter) getPendingTaskResources(environmentID string) ([]resource.ResourceConsumer, error) {
	ecsConsumers, err := c.getPendingTaskResourcesInECS(environmentID)
	if err != nil {
		return nil, err
	}

	jobConsumers, err := c.getPendingTaskResourcesInJobs(environmentID)
	if err != nil {
		return nil, err
	}

	return append(ecsConsumers, jobConsumers...), nil
}

func (c *ClusterResourceGetter) getPendingTaskResourcesInECS(environmentID string) ([]resource.ResourceConsumer, error) {
	tasks, err := c.TaskLogic.ListTasks()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(tasks); i++ {
		if tasks[i].EnvironmentID != environmentID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			i--
		}
	}

	resourceConsumers := []resource.ResourceConsumer{}
	for _, task := range tasks {
		details, err := c.TaskLogic.GetTask(task.TaskID)
		if err != nil {
			return nil, err
		}

		if details.PendingCount > 0 {
			resources, err := c.getResourcesForDeploy(details.DeployID)
			if err != nil {
				return nil, err
			}

			for i := 0; i < int(details.PendingCount); i++ {
				for _, r := range resources {
					id := fmt.Sprintf("%s.%s.%s.%d", task.TaskID, details.DeployID, r.ID, i)
					consumer := resource.NewResourceConsumer(id, r.Memory, r.Ports)
					resourceConsumers = append(resourceConsumers, consumer)
				}
			}
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
				resources, err := c.getResourcesForDeploy(req.DeployID)
				if err != nil {
					return nil, err
				}

				// note that this isn't exact if the job has started some, but not all of the tasks
				for i := 0; i < int(req.Copies); i++ {
					for _, r := range resources {
						id := fmt.Sprintf("job%s.%s.%s.%d", req.TaskName, req.DeployID, r.ID, i)
						consumer := resource.NewResourceConsumer(id, r.Memory, r.Ports)
						resourceConsumers = append(resourceConsumers, consumer)
					}
				}
			}
			// marshal request
			// get environment id
			// get deploy id
		}
	}

	return resourceConsumers, nil
}

func (c *ClusterResourceGetter) getResourcesForDeploy(deployID string) ([]resource.ResourceConsumer, error) {
	// todo: cache deploy details

	d, err := c.DeployLogic.GetDeploy(deployID)
	if err != nil {
		return nil, err
	}

	deploy, err := ecsbackend.MarshalDeploy(d.Dockerrun)
	if err != nil {
		return nil, err
	}

	consumers := make([]resource.ResourceConsumer, len(deploy.ContainerDefinitions))
	for _, container := range deploy.ContainerDefinitions {
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

		resource := resource.NewResourceConsumer(*container.Name, memory, ports)
		consumers = append(consumers, resource)

		fmt.Printf("Consumer: %#v\n", resource)
	}

	return consumers, nil
}
