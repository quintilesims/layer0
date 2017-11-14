package aws

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// Create registers an ECS Task Definition using the specified Create Deploy Request.
// The Create Deploy Request contains the name of the Deploy and the JSON
// representation of the Task Definition to create.
func (d *DeployProvider) Create(req models.CreateDeployRequest) (string, error) {
	deployID := entityIDGenerator(req.DeployName)
	familyName := addLayer0Prefix(d.Config.Instance(), req.DeployName)

	renderedTaskDefinition, err := d.renderTaskDefinition(req.DeployFile, familyName)
	if err != nil {
		return "", err
	}

	taskDefinition, err := d.createTaskDefinition(renderedTaskDefinition)
	if err != nil {
		return "", err
	}

	version := int(aws.Int64Value(taskDefinition.Revision))
	taskDefinitionARN := aws.StringValue(taskDefinition.TaskDefinitionArn)
	if err := d.createTags(deployID, req.DeployName, strconv.Itoa(version), taskDefinitionARN); err != nil {
		return "", err
	}

	return deployID, nil
}

func (d *DeployProvider) createTaskDefinition(taskDefinition *ecs.TaskDefinition) (*ecs.TaskDefinition, error) {
	input := &ecs.RegisterTaskDefinitionInput{}
	input.SetFamily(aws.StringValue(taskDefinition.Family))
	input.SetTaskRoleArn(aws.StringValue(taskDefinition.TaskRoleArn))
	input.SetContainerDefinitions(taskDefinition.ContainerDefinitions)
	input.SetVolumes(taskDefinition.Volumes)
	input.SetPlacementConstraints(taskDefinition.PlacementConstraints)
	if nm := aws.StringValue(taskDefinition.NetworkMode); nm != "" {
		input.SetNetworkMode(nm)
	}

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

	if len(taskDefinition.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("Deploy must have at least one container definition")
	}

	taskDefinition.SetFamily(familyName)
	for _, container := range taskDefinition.ContainerDefinitions {
		if container.LogConfiguration == nil {
			logConfig := &ecs.LogConfiguration{}
			logConfig.SetLogDriver("awslogs")
			logConfig.SetOptions(map[string]*string{
				"awslogs-group":         aws.String(d.Config.LogGroupName()),
				"awslogs-region":        aws.String(d.Config.Region()),
				"awslogs-stream-prefix": aws.String("l0"),
			})

			container.SetLogConfiguration(logConfig)
		}
	}

	return taskDefinition, nil
}

func (d *DeployProvider) createTags(deployID, deployName, deployVersion, taskDefinitionARN string) error {
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
			Value:      taskDefinitionARN,
		},
	}

	for _, tag := range tags {
		if err := d.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}
