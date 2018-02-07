package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// Create runs an ECS Task using the specified Create Task Request. The Create Task
// Request contains the Task name, the Environment ID, and the Deploy ID.
// The Deploy ID is used to look up the ECS Task Definition family and version of the
// Task to run.
func (t *TaskProvider) Create(req models.CreateTaskRequest) (string, error) {
	taskID := entityIDGenerator(req.TaskName)
	fqEnvironmentID := addLayer0Prefix(t.Config.Instance(), req.EnvironmentID)

	var securityGroupIDs []*string
	environmentSecurityGroupName := getEnvironmentSGName(fqEnvironmentID)
	environmentSecurityGroup, err := readSG(t.AWS.EC2, environmentSecurityGroupName)
	if err != nil {
		return "", err
	}

	securityGroupIDs = append(securityGroupIDs, environmentSecurityGroup.GroupId)

	launchType, err := getLaunchTypeFromEnvironmentID(t.TagStore, req.EnvironmentID)
	if err != nil {
		return "", err
	}

	deployName, deployVersion, err := lookupDeployNameAndVersion(t.TagStore, req.DeployID)
	if err != nil {
		return "", err
	}

	clusterName := fqEnvironmentID
	startedBy := t.Config.Instance()
	taskDefinitionFamily := addLayer0Prefix(t.Config.Instance(), deployName)
	taskDefinitionVersion := deployVersion
	taskOverrides := convertContainerOverrides(req.ContainerOverrides)

	privateSubnets := t.Config.PrivateSubnets()
	subnets := make([]*string, len(privateSubnets))
	for i := range privateSubnets {
		subnets[i] = aws.String(privateSubnets[i])
	}

	task, err := t.runTask(
		clusterName,
		launchType,
		startedBy,
		taskDefinitionFamily,
		taskDefinitionVersion,
		securityGroupIDs,
		subnets,
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
	launchType,
	startedBy,
	taskDefinitionFamily,
	taskDefinitionRevision string,
	securityGroupIDs,
	subnets []*string,
	overrides *ecs.TaskOverride,
) (*ecs.Task, error) {
	input := &ecs.RunTaskInput{}
	input.SetCluster(clusterName)
	input.SetLaunchType(launchType)
	input.SetStartedBy(startedBy)
	input.SetOverrides(overrides)

	// LAUNCH TYPE TESTING

	awsvpcConfig := &ecs.AwsVpcConfiguration{}
	awsvpcConfig.SetAssignPublicIp(ecs.AssignPublicIpDisabled) // DISABLED by default

	// environment's security group; add load balancer sg as well if exists and is public
	// (look into the security groups of a public vs private load balancer)
	awsvpcConfig.SetSecurityGroups(securityGroupIDs)

	// get from config (maybe config.privateSubnets or something)
	awsvpcConfig.SetSubnets(subnets)

	networkConfig := &ecs.NetworkConfiguration{}
	networkConfig.SetAwsvpcConfiguration(awsvpcConfig)

	input.SetNetworkConfiguration(networkConfig)

	// possibly unnecessary
	input.SetPlatformVersion("LATEST")

	// may also need to do this (unsure if these are created by default):
	// deploymentConfig := &ecs.DeploymentConfiguration{}
	// deploymentConfig.Set[somethingsabouthealthypercent]
	// input.SetDeploymentConfiguration(deploymentConfig)

	// END OF LAUNCH TYPE TESTING

	taskFamilyRevision := fmt.Sprintf("%s:%s", taskDefinitionFamily, taskDefinitionRevision)
	input.SetTaskDefinition(taskFamilyRevision)

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
