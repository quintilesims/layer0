package aws

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (d *DeployProvider) Create(req models.CreateDeployRequest) (*models.Deploy, error) {
	deployID := generateEntityID(req.DeployName)
	familyName := addLayer0Prefix(d.Config.Instance(), req.DeployName)

	taskDefinition, err := d.renderTaskDefinition(req.DeployFile, familyName)
	if err != nil {
		return nil, err
	}

	taskDefinitionOutput, err := d.createTaskDefinition(taskDefinition)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(taskDefinitionOutput)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract deploy file: %s", err.Error())
	}

	deploy := &models.Deploy{
		DeployName: req.DeployName,
		DeployID:   deployID,
		Version:    strconv.FormatInt(aws.Int64Value(taskDefinitionOutput.Revision), 10),
		DeployFile: bytes,
	}

	taskDefinitionARN := aws.StringValue(taskDefinitionOutput.TaskDefinitionArn)

	if err := d.createTags(deploy.DeployName, deploy.DeployID, deploy.Version, taskDefinitionARN); err != nil {
		return deploy, err
	}

	return deploy, nil
}

func (d *DeployProvider) createTaskDefinition(taskDefinitionRequest *ecs.TaskDefinition) (*ecs.TaskDefinition, error) {

	input := &ecs.RegisterTaskDefinitionInput{}
	input.SetFamily(aws.StringValue(taskDefinitionRequest.Family))
	input.SetTaskRoleArn(aws.StringValue(taskDefinitionRequest.TaskRoleArn))
	input.SetNetworkMode(aws.StringValue(taskDefinitionRequest.NetworkMode))
	input.SetContainerDefinitions(taskDefinitionRequest.ContainerDefinitions)
	input.SetVolumes(taskDefinitionRequest.Volumes)
	input.SetPlacementConstraints(taskDefinitionRequest.PlacementConstraints)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := d.AWS.ECS.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return output.TaskDefinition, nil
}

func (d *DeployProvider) renderTaskDefinition(body []byte, familyName string) (*ecs.TaskDefinition, error) {
	var taskDefinition *ecs.TaskDefinition

	if err := json.Unmarshal(body, &taskDefinition); err != nil {
		return nil, fmt.Errorf("Failed to decode deploy: %s", err.Error())
	}

	taskDefinition.SetFamily(familyName)

	if len(taskDefinition.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("Deploy must have at least one container definition")
	}

	for _, container := range taskDefinition.ContainerDefinitions {
		if container.LogConfiguration == nil {
			logConfig := &ecs.LogConfiguration{
				LogDriver: aws.String("awslogs"),
				Options: map[string]*string{
					"awslogs-group":         aws.String(fmt.Sprintf("l0-%s", d.Config.Instance())),
					"awslogs-region":        aws.String(d.Config.Region()),
					"awslogs-stream-prefix": aws.String("l0"),
				},
			}
			container.SetLogConfiguration(logConfig)
		}
	}

	return taskDefinition, nil
}

func (d *DeployProvider) createTags(deployName, deployID, deployVersion, taskDefinitionArn string) error {
	tags := []models.Tag{
		{
			EntityID:   deployID,
			EntityType: "deploy",
			Key:        "name",
			Value:      deployName,
		},
		{
			EntityID:   deployID,
			EntityType: "deploy",
			Key:        "version",
			Value:      deployVersion,
		},
		{
			EntityID:   deployID,
			EntityType: "deploy",
			Key:        "arn",
			Value:      taskDefinitionArn,
		},
	}

	for _, tag := range tags {
		if err := d.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}
