package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

// Create runs an ECS Task using the specified CreateTaskRequest.
// The CreateTaskRequest contains the ContainerOverrides, the DeployID, the
// EnvironmentID, the TaskName, and the TaskType.
//
// The Deploy ID is used to look up the ECS TaskDefinition Family and Version
// of the Task to run. The TaskType parameter is one of "stateful" or "stateless"
// and indicates which ECS LaunchType the user wishes to use ("EC2" or "FARGATE"
// respectively).
//
// Create does not generate any custom errors of its own, but will bubble up errors
// found in its helper functions as well as errors returned by AWS.
func (t *TaskProvider) Create(req models.CreateTaskRequest) (string, error) {
	taskID := entityIDGenerator(req.TaskName)
	fqEnvironmentID := addLayer0Prefix(t.Config.Instance(), req.EnvironmentID)
	clusterName := fqEnvironmentID
	startedBy := t.Config.Instance()
	taskOverrides := convertContainerOverrides(req.ContainerOverrides)

	taskDefinitionARN, err := lookupTaskDefinitionARNFromDeployID(t.TagStore, req.DeployID)
	if err != nil {
		return "", err
	}

	taskDefinition, err := describeTaskDefinition(t.AWS.ECS, taskDefinitionARN)
	if err != nil {
		return "", err
	}

	networkMode := aws.StringValue(taskDefinition.NetworkMode)

	var securityGroupIDs []*string
	var subnets []string
	if networkMode == ecs.NetworkModeAwsvpc {
		environmentSecurityGroupName := getEnvironmentSGName(fqEnvironmentID)
		environmentSecurityGroup, err := readSG(t.AWS.EC2, environmentSecurityGroupName)
		if err != nil {
			return "", err
		}

		securityGroupIDs = append(securityGroupIDs, environmentSecurityGroup.GroupId)

		subnets = t.Config.PrivateSubnets()
	}

	task, err := t.runTask(
		clusterName,
		startedBy,
		taskDefinitionARN,
		networkMode,
		req.Stateful,
		subnets,
		securityGroupIDs,
		taskOverrides)
	if err != nil {
		return "", err
	}

	taskARN := aws.StringValue(task.TaskArn)
	if err := t.createTags(taskID, req.TaskName, req.EnvironmentID, taskARN); err != nil {
		return "", err
	}

	return taskID, nil
}

func convertContainerOverrides(overrides []models.ContainerOverride) *ecs.TaskOverride {
	ecsOverrides := make([]*ecs.ContainerOverride, len(overrides))
	for i, o := range overrides {
		environment := []*ecs.KeyValuePair{}
		for name, value := range o.EnvironmentOverrides {
			kvp := &ecs.KeyValuePair{}
			kvp.SetName(name)
			kvp.SetValue(value)

			environment = append(environment, kvp)
		}

		ecsOverride := &ecs.ContainerOverride{}
		ecsOverride.SetName(o.ContainerName)
		ecsOverride.SetEnvironment(environment)

		ecsOverrides[i] = ecsOverride
	}

	taskOverride := &ecs.TaskOverride{}
	taskOverride.SetContainerOverrides(ecsOverrides)

	return taskOverride
}

func (t *TaskProvider) runTask(
	clusterName,
	startedBy,
	taskDefinitionARN,
	networkMode string,
	stateful bool,
	subnets []string,
	securityGroupIDs []*string,
	overrides *ecs.TaskOverride,
) (*ecs.Task, error) {
	input := &ecs.RunTaskInput{}
	input.SetCluster(clusterName)
	input.SetStartedBy(startedBy)
	input.SetOverrides(overrides)

	launchType := ecs.LaunchTypeEc2
	if !stateful {
		launchType = ecs.LaunchTypeFargate
		input.SetPlatformVersion(config.DefaultFargatePlatformVersion)
	}

	input.SetLaunchType(launchType)

	if networkMode == ecs.NetworkModeAwsvpc {
		s := make([]*string, len(subnets))
		for i := range subnets {
			s[i] = aws.String(subnets[i])
		}

		awsvpcConfig := &ecs.AwsVpcConfiguration{}
		awsvpcConfig.SetAssignPublicIp(ecs.AssignPublicIpDisabled)
		awsvpcConfig.SetSecurityGroups(securityGroupIDs)
		awsvpcConfig.SetSubnets(s)

		networkConfig := &ecs.NetworkConfiguration{}
		networkConfig.SetAwsvpcConfiguration(awsvpcConfig)

		input.SetNetworkConfiguration(networkConfig)
	}

	input.SetTaskDefinition(taskDefinitionARN)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := t.AWS.ECS.RunTask(input)
	if err != nil {
		return nil, err
	}

	if len(output.Failures) > 0 {
		return nil, fmt.Errorf("Failed to create task: %s", aws.StringValue(output.Failures[0].Reason))
	}

	return output.Tasks[0], nil
}

func (t *TaskProvider) createTags(taskID, taskName, environmentID, taskARN string) error {
	tags := []models.Tag{
		{
			EntityID:   taskID,
			EntityType: "task",
			Key:        "name",
			Value:      taskName,
		},
		{
			EntityID:   taskID,
			EntityType: "task",
			Key:        "environment_id",
			Value:      environmentID,
		},
		{
			EntityID:   taskID,
			EntityType: "task",
			Key:        "arn",
			Value:      taskARN,
		},
	}

	for _, tag := range tags {
		if err := t.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}
