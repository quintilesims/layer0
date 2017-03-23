package ecsbackend

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/api/scheduler/resource"
	"github.com/quintilesims/layer0/common/aws/autoscaling"
	"github.com/quintilesims/layer0/common/aws/ec2"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/logutils"
	"github.com/zpatrick/go-bytesize"
	"strconv"
)

type ECSResourceManager struct {
	ECS         ecs.Provider
	Autoscaling autoscaling.Provider
	logger      *logrus.Logger
}

func NewECSResourceManager(e ecs.Provider, a autoscaling.Provider) *ECSResourceManager {
	return &ECSResourceManager{
		ECS:         e,
		Autoscaling: a,
		logger:      logutils.NewStandardLogger("ECS Resource Manager"),
	}
}

func (r *ECSResourceManager) GetProviders(environmentID string) ([]*resource.ResourceProvider, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()

	instanceARNs, err := r.ECS.ListContainerInstances(ecsEnvironmentID.String())
	if err != nil {
		return nil, err
	}

	resourceProviders := []*resource.ResourceProvider{}
	if len(instanceARNs) > 0 {
		instances, err := r.ECS.DescribeContainerInstances(ecsEnvironmentID.String(), instanceARNs)
		if err != nil {
			return nil, err
		}

		for _, instance := range instances {
			provider, ok := r.getResourceProvider(ecsEnvironmentID, instance)
			if ok {
				resourceProviders = append(resourceProviders, provider)
			}
		}
	}

	return resourceProviders, nil
}

func (r *ECSResourceManager) getResourceProvider(ecsEnvironmentID id.ECSEnvironmentID, instance *ecs.ContainerInstance) (*resource.ResourceProvider, bool) {
	instanceID := pstring(instance.Ec2InstanceId)
	if pstring(instance.Status) != "ACTIVE" {
		r.logger.Warnf("Instance %d is not in the 'ACTIVE' state", instanceID)
		return nil, false
	}

	if !pbool(instance.AgentConnected) {
		r.logger.Errorf("Instance %d agent is disconnected")
		return nil, false
	}

	// this is non-intuitive, but the ports being used by tasks are kept in
	// instance.ReminaingResources, not instance.RegisteredResources
	var usedPorts []int
	var availableMemory bytesize.Bytesize
	for _, resource := range instance.RemainingResources {
		switch pstring(resource.Name) {
		case "MEMORY":
			v := pint64(resource.IntegerValue)
			availableMemory = bytesize.MiB * bytesize.Bytesize(v)

		case "PORTS":
			for _, p := range resource.StringSetValue {
				port, err := strconv.Atoi(pstring(p))
				if err != nil {
					r.logger.Errorf("Instance %d: Failed to convert port to int: %v\n", instanceID, err)
					continue
				}

				usedPorts = append(usedPorts, port)
			}
		}
	}

	inUse := pint64(instance.PendingTasksCount)+pint64(instance.RunningTasksCount) > 0
	provider := resource.NewResourceProvider(instanceID, inUse, availableMemory, usedPorts)

	r.logger.Debugf("Environment '%s' generated provider: %#v\n", ecsEnvironmentID, provider)
	return provider, true
}

func (r *ECSResourceManager) CalculateNewProvider(environmentID string) (*resource.ResourceProvider, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	group, err := r.Autoscaling.DescribeAutoScalingGroup(ecsEnvironmentID.String())
	if err != nil {
		return nil, err
	}

	config, err := r.Autoscaling.DescribeLaunchConfiguration(*group.LaunchConfigurationName)
	if err != nil {
		return nil, err
	}

	memory, ok := ec2.InstanceSizes[*config.InstanceType]
	if !ok {
		return nil, fmt.Errorf("Environment %s is using unknown instance type '%s'", *config.InstanceType)
	}

	// these ports are automatically used by the ecs agent
	defaultPorts := []int{
		22,
		2376,
		2375,
		51678,
		51679,
	}

	return resource.NewResourceProvider("<new instance>", false, memory, defaultPorts), nil
}

func (r *ECSResourceManager) ScaleTo(environmentID string, scale int, unusedProviders []*resource.ResourceProvider) (int, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	asg, err := r.Autoscaling.DescribeAutoScalingGroup(ecsEnvironmentID.String())
	if err != nil {
		return 0, err
	}

	currentCapacity := int(pint64(asg.DesiredCapacity))

	switch {
	case scale > currentCapacity:
		r.logger.Debugf("Environment %s is attempting to scale up to size %d", ecsEnvironmentID, scale)
		return r.scaleUp(ecsEnvironmentID, scale, asg)
	case scale < currentCapacity:
		r.logger.Debugf("Environment %s is attempting to scale down to size %d", ecsEnvironmentID, scale)
		return r.scaleDown(ecsEnvironmentID, scale, asg, unusedProviders)
	default:
		r.logger.Debugf("Environment %s is at desired scale of %d. No scaling action required.", ecsEnvironmentID, scale)
		return currentCapacity, nil
	}
}

func (r *ECSResourceManager) scaleUp(ecsEnvironmentID id.ECSEnvironmentID, scale int, asg *autoscaling.Group) (int, error) {
	maxCapacity := int(pint64(asg.MaxSize))
	if scale > maxCapacity {
		if err := r.Autoscaling.UpdateAutoScalingGroupMaxSize(*asg.AutoScalingGroupName, scale); err != nil {
			return 0, err
		}
	}

	if err := r.Autoscaling.SetDesiredCapacity(*asg.AutoScalingGroupName, scale); err != nil {
		return 0, err
	}

	return scale, nil
}

func (r *ECSResourceManager) scaleDown(ecsEnvironmentID id.ECSEnvironmentID, scale int, asg *autoscaling.Group, unusedProviders []*resource.ResourceProvider) (int, error) {
	minCapacity := int(pint64(asg.MinSize))
	if scale < minCapacity {
		r.logger.Warnf("Scale %d is below the minimum capacity of %d. Setting desired capacity to %d.", scale, minCapacity, minCapacity)
		scale = minCapacity
	}

	currentCapacity := int(pint64(asg.DesiredCapacity))
	if scale == currentCapacity {
		r.logger.Debugf("Environment %s is at desired scale of %d. No scaling action required.", ecsEnvironmentID, scale)
		return scale, nil
	}

	if scale < currentCapacity {
		if err := r.Autoscaling.SetDesiredCapacity(*asg.AutoScalingGroupName, scale); err != nil {
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
		r.logger.Debugf("Environment %s terminating unused instance '%s'", ecsEnvironmentID, unusedProvider.ID)

		if _, err := r.Autoscaling.TerminateInstanceInAutoScalingGroup(unusedProvider.ID, false); err != nil {
			return 0, err
		}
	}

	return scale, nil
}
