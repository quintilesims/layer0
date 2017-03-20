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

func (p *ECSResourceManager) GetProviders(environmentID string) ([]*resource.ResourceProvider, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()

	instanceARNs, err := p.ECS.ListContainerInstances(ecsEnvironmentID.String())
	if err != nil {
		return nil, err
	}

	resourceProviders := []*resource.ResourceProvider{}
	if len(instanceARNs) > 0 {
		instances, err := p.ECS.DescribeContainerInstances(ecsEnvironmentID.String(), instanceARNs)
		if err != nil {
			return nil, err
		}

		for _, instance := range instances {
			provider, ok := p.getResourceProvider(ecsEnvironmentID, instance)
			if ok {
				resourceProviders = append(resourceProviders, provider)
			}
		}
	}

	return resourceProviders, nil
}

func (c *ECSResourceManager) getResourceProvider(ecsEnvironmentID id.ECSEnvironmentID, instance *ecs.ContainerInstance) (*resource.ResourceProvider, bool) {
	instanceID := pstring(instance.Ec2InstanceId)
	if pstring(instance.Status) != "ACTIVE" {
		c.logger.Warnf("Instance %d is not in the 'ACTIVE' state", instanceID)
		return nil, false
	}

	// todo: what if instance.AgentConnected == false?
	// declare it as a not in-use instance?
	// do nothing, capture this in scale function?

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
					c.logger.Errorf("Instance %d: Failed to convert port to int: %v\n", instanceID, err)
					continue
				}

				usedPorts = append(usedPorts, port)
			}
		}
	}

	inUse := pint64(instance.PendingTasksCount)+pint64(instance.RunningTasksCount) > 0
	provider := resource.NewResourceProvider(instanceID, inUse, availableMemory, usedPorts)

	c.logger.Debugf("Environment '%s' generated provider: %#v\n", ecsEnvironmentID, provider)
	return provider, true
}

func (c *ECSResourceManager) AddNewProvider(environmentID string) (*resource.ResourceProvider, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	group, err := c.Autoscaling.DescribeAutoScalingGroup(ecsEnvironmentID.String())
	if err != nil {
		return nil, err
	}

	config, err := c.Autoscaling.DescribeLaunchConfiguration(*group.LaunchConfigurationName)
	if err != nil {
		return nil, err
	}

	memory, ok := ec2.InstanceSizes[*config.InstanceType]
	if !ok {
		return nil, fmt.Errorf("Environment %s is using unknown instance type '%s'", *config.InstanceType)
	}

	// todo: add default ports
	return resource.NewResourceProvider("<new instance>", false, memory, nil), nil
}

func (c *ECSResourceManager) ScaleTo(environmentID string, scale int, unusedProviders []*resource.ResourceProvider) (int, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	asg, err := c.Autoscaling.DescribeAutoScalingGroup(ecsEnvironmentID.String())
	if err != nil {
		return 0, err
	}

	currentCapacity := int(pint64(asg.DesiredCapacity))

	switch {
	case scale > currentCapacity:
		c.logger.Debugf("Environment %s is attempting to scale up to size %d", ecsEnvironmentID, scale)
		return c.scaleUp(ecsEnvironmentID, scale, asg)
	case scale < currentCapacity:
		c.logger.Debugf("Environment %s is attempting to scale down to size %d", ecsEnvironmentID, scale)
		return c.scaleDown(ecsEnvironmentID, scale, asg, unusedProviders)
	default:
		c.logger.Debugf("Environment %s is at desired scale of %d. No scaling action required.", ecsEnvironmentID, scale)
		return currentCapacity, nil
	}
}

func (c *ECSResourceManager) scaleUp(ecsEnvironmentID id.ECSEnvironmentID, scale int, asg *autoscaling.Group) (int, error) {
	maxCapacity := int(pint64(asg.MaxSize))
	if scale > maxCapacity {
		if err := c.Autoscaling.UpdateAutoScalingGroupMaxSize(*asg.AutoScalingGroupName, scale); err != nil {
			return 0, err
		}
	}

	if err := c.Autoscaling.SetDesiredCapacity(asg.AutoScalingGroupName, scale); err != nil {
		return 0, err
	}

	return scale, nil
}

func (c *ECSResourceManager) scaleDown(ecsEnvironmentID id.ECSEnvironmentID, scale int, asg *autoscaling.Group, unusedProviders []*resource.ResourceProvider) (int, error) {
	minCapacity := int(pint64(asg.MinSize))
	if scale < minCapacity {
		c.logger.Warnf("Scale %d is below the minimum capacity of %d. Setting desired capacity to %d.", scale, minCapacity, minCapacity)
		scale = minCapacity
	}

	currentCapacity := int(pint64(asg.DesiredCapacity))
	if scale == currentCapacity {
		c.logger.Debugf("Environment %s is at desired scale of %d. No scaling action required.", ecsEnvironmentID, scale)
		return scale, nil
	}

	if scale < currentCapacity {
		if err := c.Autoscaling.SetDesiredCapacity(asg.AutoScalingGroupName, scale); err != nil {
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
		c.logger.Debugf("Environment %s terminating unused instance '%s'", ecsEnvironmentID, unusedProvider.ID)

		if _, err := c.Autoscaling.TerminateInstanceInAutoScalingGroup(unusedProvider.ID, false); err != nil {
			return 0, err
		}
	}

	return scale, nil
}
