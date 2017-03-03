package ecsbackend

import (
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/autoscaling"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/models"
	"strconv"
)

type ECSClusterScaler struct {
	ECS         ecs.Provider
	Autoscaling autoscaling.Provider
	Backend     backend.Backend
}

type ClusterScaler interface {
	TriggerScalingAlgorithm(ecsEnvironmentID id.ECSEnvironmentID, ecsDeployID *id.ECSDeployID, count int) (int, bool, error)
}

func NewECSClusterScaler(ecs ecs.Provider, autoscaling autoscaling.Provider, b backend.Backend) *ECSClusterScaler {
	return &ECSClusterScaler{
		ECS:         ecs,
		Autoscaling: autoscaling,
		Backend:     b,
	}
}

type Resources struct {
	CPU, Memory     int64
	Ports, UdpPorts []*string
}

type InstanceResources struct {
	Resources
	TaskCount int64
}

func (this *ECSClusterScaler) TriggerScalingAlgorithm(ecsEnvironmentID id.ECSEnvironmentID, ecsDeployID *id.ECSDeployID, count int) (int, bool, error) {
	log.Debugf("Running scaling algorithm for %s, %v, %d", ecsEnvironmentID, ecsDeployID, count)

	resourcesAvail, err := this.getCurrentResources(ecsEnvironmentID)
	if err != nil {
		return 0, false, err
	}

	newInstancesNeeded, hasUnallocatedTasks, err := this.calculateNewInstancesNeeded(ecsEnvironmentID, ecsDeployID, count, resourcesAvail)
	if err != nil {
		return 0, false, err
	}

	// if we don't have space for (task+a task of HEADROOM) size, add the required
	// # of CIs
	newInstancesAdded := 0
	if newInstancesNeeded > 0 {
		asg, err := this.Autoscaling.DescribeAutoScalingGroup(ecsEnvironmentID.String())
		if err != nil {
			return 0, false, err
		}

		// we may just be waiting for instances to come up, so don't blindly add
		// capacity to the ASG
		totalInstancesNeeded := newInstancesNeeded + len(resourcesAvail)
		if totalInstancesNeeded > int(*asg.DesiredCapacity) {
			clusterName := ecsEnvironmentID.String()

			if err := this.Autoscaling.UpdateAutoScalingGroupMaxSize(clusterName, totalInstancesNeeded); err != nil {
				return 0, false, err
			}

			newInstancesAdded = totalInstancesNeeded - int(*asg.DesiredCapacity)
		}
	}

	return newInstancesAdded, hasUnallocatedTasks, err
}

func (this *ECSClusterScaler) getCurrentResources(ecsEnvironmentID id.ECSEnvironmentID) ([]*InstanceResources, error) {
	// get list of CIs for cluster
	instances, err := this.ECS.ListContainerInstances(ecsEnvironmentID.String())
	if err != nil {
		return nil, err
	}

	resourcesAvail := []*InstanceResources{}
	if len(instances) > 0 {
		// for each container instance, see if it has space for the
		// task (cpu, mem, tcp ports, udp ports, agent connected)
		instanceDescs, err := this.ECS.DescribeContainerInstances(ecsEnvironmentID.String(), instances)
		if err != nil {
			if len(instanceDescs) == 0 {
				return nil, err
			} else {
				log.Errorf("[ClusteScaler] Some, but not all container instances failed to describe.  Ignoring the failed instances: %v", err)
			}
		}

		for _, instance := range instanceDescs {
			if *instance.AgentConnected && *instance.Status == "ACTIVE" {
				newResource := &InstanceResources{}
				newResource.TaskCount = *instance.PendingTasksCount + *instance.RunningTasksCount

				for _, resource := range instance.RemainingResources {
					switch *resource.Name {
					case "CPU":
						newResource.CPU = *resource.IntegerValue
					case "MEMORY":
						newResource.Memory = *resource.IntegerValue
					case "PORTS":
						newResource.Ports = resource.StringSetValue
					case "PORTS_UDP":
						newResource.UdpPorts = resource.StringSetValue
					}
				}

				resourcesAvail = append(resourcesAvail, newResource)
			}
		}
	}

	return resourcesAvail, nil
}

// Provides a measure of the number of new instances needed to run `count` copies of `task`
// Returns (newInstancesNeeded, hasUnallocatedTasks, error)
// newInstancesNeeded - Number of instances to add to fit `count` copies of `familyAndRevision` in `clusterName`
// hasUnallocatedTasks - Boolean whether the cluster has unallocated tasks
// (tasks that are neither RUNNING, STOPPED, or PENDING, but rather in a queue to be started).
//  A good signal to not scale the cluster down
func (this *ECSClusterScaler) calculateNewInstancesNeeded(
	ecsEnvironmentID id.ECSEnvironmentID,
	ecsDeployID *id.ECSDeployID,
	count int,
	availableInstanceResources []*InstanceResources,
) (int, bool, error) {

	// Account for pending tasks
	serviceARNs, err := this.ECS.ListServices(ecsEnvironmentID.String())
	if err != nil {
		return 0, false, err
	}

	newInstancesRequired := 0
	hasUnallocatedTasks := false
	if len(serviceARNs) > 0 {
		services, err := this.ECS.DescribeServices(ecsEnvironmentID.String(), serviceARNs)
		if err != nil {
			return 0, false, err
		}

		type PendingTask struct {
			TaskDefinition *string
			Count          int64
		}

		pendingTasks := []*PendingTask{}

		// Inspect deployments for services for diff between running+pending == desired counts
		for _, service := range services {
			for _, deploy := range service.Deployments {

				// for PRIMARY that aren't fully running make sure there's space in the cluster for them
				if *deploy.Status == "PRIMARY" {
					if *deploy.DesiredCount != *deploy.PendingCount+*deploy.RunningCount {
						hasUnallocatedTasks = true

						// still allocating, make sure there's space
						requiredSpace := *deploy.DesiredCount - (*deploy.PendingCount + *deploy.RunningCount)
						pendingTasks = append(pendingTasks, &PendingTask{deploy.TaskDefinition, requiredSpace})
					}
				}
			}
		}

		for _, task := range pendingTasks {
			ecsTaskDeployID := id.ECSDeployID(*task.TaskDefinition)
			pendingTaskResources, err := this.calculateTaskResources(ecsTaskDeployID)
			if err != nil {
				return 0, false, err
			}

			for i := int64(0); i < task.Count; i++ {
				if this.taskNeedsNewInstance(availableInstanceResources, pendingTaskResources) {
					newInstancesRequired += 1
				}
			}
		}
	}

	// Account for Layer0 1-Off Tasks resource commit
	// Use the backend call which will include all running and pending (including
	// unallocated) tasks
	tasks, err := this.getAllTasks()
	if err != nil {
		return 0, false, err
	}

	for _, task := range tasks {
		ecsTaskDeployID := id.L0DeployID(task.DeployID).ECSDeployID()
		pendingTaskResources, err := this.calculateTaskResources(ecsTaskDeployID)
		if err != nil {
			return 0, false, err
		}

		// only claim resources for running or about to run
		// desired would potentially include stopped tasks
		for _, c := range task.Copies {
			if c.Reason == ClusterCapacityReason {
				hasUnallocatedTasks = true
				if this.taskNeedsNewInstance(availableInstanceResources, pendingTaskResources) {
					newInstancesRequired += 1
				}
			}

			for _, detail := range c.Details {
				// since this variable doubles as "hasPendingTasks" in right_sizer,
				// we mark it true here
				if detail.LastStatus == "PENDING" {
					hasUnallocatedTasks = true
				}
			}
		}
	}

	// Finally, see if the new task we want to run fits
	if ecsDeployID != nil {
		taskResources, err := this.calculateTaskResources(*ecsDeployID)
		if err != nil {
			return 0, false, err
		}

		for i := 0; i < count; i++ {
			if this.taskNeedsNewInstance(availableInstanceResources, taskResources) {
				newInstancesRequired += 1
			}
		}
	}

	return newInstancesRequired, hasUnallocatedTasks, nil
}

func (this *ECSClusterScaler) calculateTaskResources(ecsDeployID id.ECSDeployID) (*Resources, error) {
	taskResources := &Resources{}

	// lookup the taskdef and get the requirements from it
	taskDef, err := this.ECS.DescribeTaskDefinition(ecsDeployID.TaskDefinition())
	if err != nil {
		return nil, err
	}

	// accumulate the total resources needed for the containers in this task definition
	for _, container := range taskDef.ContainerDefinitions {
		taskResources.CPU += *container.Cpu

		var memory int64
		if container.MemoryReservation != nil {
			memory = *container.MemoryReservation
		}

<<<<<<< HEAD
		// override if container.Memory is specified because that is a hard limit 
=======
		// override if container.Memory is specified because that is a hard limit
>>>>>>> 579808dcf34d524874fffa1b885e18866ebddae9
		if container.Memory != nil {
			memory = *container.Memory
		}

		taskResources.Memory += memory

		for _, mapping := range container.PortMappings {
			if mapping.HostPort != nil {
				if mapping.Protocol == nil || *mapping.Protocol == "tcp" {
					hostPort := strconv.FormatInt(
						*mapping.HostPort, 10)
					taskResources.Ports = append(
						taskResources.Ports, &hostPort)
				} else {
					hostPort := strconv.FormatInt(
						*mapping.HostPort, 10)
					taskResources.UdpPorts = append(
						taskResources.UdpPorts, &hostPort)
				}
			}
		}
	}

	return taskResources, nil
}

// Takes the resources for `taskResources` if available, returns boolean true if it doesn't fit
// `availableInstanceResources` is mutated to no longer claim the used resource
// TODO: There is a subtle bug here that the resource removed may not actually be the resource that gets selected by the ECS scheduler
// We don't get to control CI selection for ECS services (only tasks, even then it's a sort of odd thing to do under our current scheduling strategy)
// Example:
//   available: [ 3000M, 1500M ]
//   Task1: 1500M, Task2 3000M
//   schedule task1, avail [ 3000M ], schedule task2 avail []
// ECS scheduler could spin up task1 on the box with more resource, starving task2 indefinitely
// In fact ECS scheduler will prefer to schedule onto the box with fewer running tasks (ie the one with more avail MEM)
// Even with that, does this demand a RightSizer async job to catch these stuck jobs? I think so
func (this *ECSClusterScaler) taskNeedsNewInstance(availableInstanceResources []*InstanceResources, taskResources *Resources) bool {

	// TODO check agentconnected on the instance
	// TODO check if instance has more than 50 ports mapped (appears to be an ECS limitaiton)

	usableResources := []*InstanceResources{}

OuterLoop:
	for _, resource := range availableInstanceResources {
		if taskResources.CPU <= resource.CPU {
			if taskResources.Memory <= resource.Memory {
				for _, port := range taskResources.Ports {
					for _, usedPort := range resource.Ports {
						if *usedPort == *port {
							continue OuterLoop
						}
					}
				}
				for _, port := range taskResources.UdpPorts {
					for _, usedPort := range resource.UdpPorts {
						if *usedPort == *port {
							continue OuterLoop
						}
					}
				}
				// looks like we can fit this resource
				usableResources = append(usableResources, resource)

				// could update this resource with it's new commit and see if we can fit another...for now let's just make new machines
			}
		}
	}

	// This is actually just a guess as to which resource is going to get used based on the notes above.  We're counting on RightSizer to fix this in the worst case
	var resourceWithMinTasks *InstanceResources
	for _, resource := range usableResources {
		if resourceWithMinTasks == nil || resourceWithMinTasks.TaskCount > resource.TaskCount {
			resourceWithMinTasks = resource
		}
	}

	if resourceWithMinTasks != nil {
		resourceWithMinTasks.CPU -= taskResources.CPU
		resourceWithMinTasks.Memory -= taskResources.Memory
		resourceWithMinTasks.TaskCount += 1
		resourceWithMinTasks.Ports = append(resourceWithMinTasks.Ports, taskResources.Ports...)
		resourceWithMinTasks.UdpPorts = append(resourceWithMinTasks.UdpPorts, taskResources.UdpPorts...)
	}

	return resourceWithMinTasks == nil
}

func (this *ECSClusterScaler) getAllTasks() ([]*models.Task, error) {
	summaries, err := this.Backend.ListTasks()
	if err != nil {
		return nil, err
	}

	tasks := make([]*models.Task, len(summaries))
	for i, summary := range summaries {
		task, err := this.Backend.GetTask(summary.EnvironmentID, summary.TaskID)
		if err != nil {
			return nil, err
		}

		tasks[i] = task
	}

	return tasks, nil
}
