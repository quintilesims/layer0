package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func (e *EnvironmentScaler) getActiveContainerInstanceARNsForCluster(clusterName string) ([]*string, error) {
	containerInstanceARNs := []*string{}

	listContainerInstancesPagesFN := func(output *ecs.ListContainerInstancesOutput, lastPage bool) bool {
		containerInstanceARNs = append(containerInstanceARNs, output.ContainerInstanceArns...)

		return !lastPage
	}

	listContainerInstancesInput := &ecs.ListContainerInstancesInput{}
	listContainerInstancesInput.SetCluster(clusterName)
	listContainerInstancesInput.SetStatus(ecs.ContainerInstanceStatusActive)
	if err := e.Client.ECS.ListContainerInstancesPages(listContainerInstancesInput, listContainerInstancesPagesFN); err != nil {
		return nil, err
	}

	return containerInstanceARNs, nil
}

func (e *EnvironmentScaler) getAutoScalingGroupForCluster(clusterName string) (*autoscaling.Group, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	input.SetAutoScalingGroupNames([]*string{&clusterName})

	asgs, err := e.Client.AutoScaling.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil, err
	}

	return asgs.AutoScalingGroups[0], nil
}

func (e *EnvironmentScaler) getContainerInstancesFromARNs(clusterName string, containerInstanceARNs []*string) ([]*ecs.ContainerInstance, error) {
	describeContainerInstancesInput := &ecs.DescribeContainerInstancesInput{}
	describeContainerInstancesInput.SetCluster(clusterName)
	describeContainerInstancesInput.SetContainerInstances(containerInstanceARNs)

	output, err := e.Client.ECS.DescribeContainerInstances(describeContainerInstancesInput)
	if err != nil {
		return nil, err
	}

	return output.ContainerInstances, nil
}

func (e *EnvironmentScaler) getServiceARNsForCluster(clusterName string) ([]*string, error) {
	var serviceARNs []*string

	listServicesPagesFN := func(output *ecs.ListServicesOutput, lastPage bool) bool {
		serviceARNs = append(serviceARNs, output.ServiceArns...)

		return !lastPage
	}

	listServicesInput := &ecs.ListServicesInput{}
	listServicesInput.SetCluster(clusterName)
	if err := e.Client.ECS.ListServicesPages(listServicesInput, listServicesPagesFN); err != nil {
		return nil, err
	}

	return serviceARNs, nil
}

func (e *EnvironmentScaler) getServicesFromServiceARNs(clusterName string, serviceARNs []*string) ([]*ecs.Service, error) {
	services := []*ecs.Service{}
	if len(serviceARNs) > 0 {
		// The SDK states that you can specify up to 10 services in one DescribeServices operation:
		// https://github.com/aws/aws-sdk-go/blob/v1.12.19/service/ecs/api.go#L5420
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

	return services, nil
}

func (e *EnvironmentScaler) getTaskARNsForCluster(clusterName, startedBy, status string) ([]string, error) {
	var taskARNs []string

	fn := func(output *ecs.ListTasksOutput, lastPage bool) bool {
		for _, taskARN := range output.TaskArns {
			taskARNs = append(taskARNs, aws.StringValue(taskARN))
		}

		return !lastPage
	}

	input := &ecs.ListTasksInput{}
	input.SetCluster(clusterName)
	input.SetDesiredStatus(status)
	input.SetStartedBy(startedBy)
	if err := e.Client.ECS.ListTasksPages(input, fn); err != nil {
		return nil, err
	}

	return taskARNs, nil
}

func (e *EnvironmentScaler) getTaskDefinitionFromDeployID(deployID string) (*ecs.TaskDefinition, error) {
	input := &ecs.DescribeTaskDefinitionInput{}
	input.SetTaskDefinition(deployID)
	output, err := e.Client.ECS.DescribeTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return output.TaskDefinition, nil
}

func (e *EnvironmentScaler) setDesiredCapacityForAutoScalingGroup(asgName string, desiredScale int) error {
	input := &autoscaling.SetDesiredCapacityInput{}
	input.SetAutoScalingGroupName(asgName)
	input.SetDesiredCapacity(int64(desiredScale))

	if _, err := e.Client.AutoScaling.SetDesiredCapacity(input); err != nil {
		return err
	}

	return nil
}
func (e *EnvironmentScaler) terminateInstanceInAutoScalingGroup(instanceID string) error {
	input := &autoscaling.TerminateInstanceInAutoScalingGroupInput{}
	input.SetInstanceId(instanceID)
	input.SetShouldDecrementDesiredCapacity(false)

	if _, err := e.Client.AutoScaling.TerminateInstanceInAutoScalingGroup(input); err != nil {
		return err
	}

	return nil
}
