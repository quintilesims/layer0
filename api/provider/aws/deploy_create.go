package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (d *DeployProvider) Create(req models.CreateDeployRequest) (*models.Deploy, error) {
	// TODO: see Line 59. Which entity ID generation pattern do we prefer? ARN to ID or the
	// example below?
	deployID := generateEntityID(req.DeployName)

	if err := req.Validate(); err != nil {
		return nil, err
	}

	deployFile, err := d.createRenderedDeploy(req.DeployFile)
	if err != nil {
		return nil, err
	}

	familyName := L0DeployID(d.Config.Instance(), req.DeployName)
	if deployFile.Family != nil && deployFile.Family != aws.String(familyName) {
		return nil, fmt.Errorf("Custom family names are currently unsupported in Layer0")
	}

	input := &ecs.RegisterTaskDefinitionInput{}
	input.SetFamily(familyName)
	input.SetTaskRoleArn(aws.StringValue(deployFile.TaskRoleArn))
	input.SetNetworkMode(aws.StringValue(deployFile.NetworkMode))
	input.SetContainerDefinitions(deployFile.ContainerDefinitions)
	input.SetVolumes(deployFile.Volumes)
	input.SetPlacementConstraints(deployFile.PlacementConstraints)

	taskDefinition := d.createTaskDefinition(input)

	bytes, err := json.Marshal(taskDefinition)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract deploy file: %s", err.Error())
	}

	deploy := &models.Deploy{
		DeployID:   L0DeployID(d.Config.Instance(), deployID),
		Version:    GetRevision(d.Config.Instance(), deployID),
		DeployFile: bytes,
	}

	if err := d.TagStore.Insert(models.Tag{EntityID: deploy.DeployID, EntityType: "deploy", Key: "name", Value: req.DeployName}); err != nil {
		return deploy, err
	}

	if err := d.TagStore.Insert(models.Tag{EntityID: deploy.DeployID, EntityType: "deploy", Key: "version", Value: deploy.Version}); err != nil {
		return deploy, err
	}

	return deploy, nil
}

func (d *DeployProvider) createTaskDefinition(input *ecs.RegisterTaskDefinitionInput) *ecs.TaskDefinition {
	// TODO: use this for generating ecsDeployID or `generateEntityID` ?
	// ecsDeployID := TaskDefinitionARNToECSDeployID(*input.TaskRoleArn)

	taskDefinition := &ecs.TaskDefinition{}
	taskDefinition.SetContainerDefinitions(input.ContainerDefinitions)
	taskDefinition.SetVolumes(input.Volumes)
	taskDefinition.SetFamily(aws.StringValue(input.Family))
	taskDefinition.SetNetworkMode(aws.StringValue(input.NetworkMode))
	taskDefinition.SetTaskRoleArn(aws.StringValue(input.TaskRoleArn))
	taskDefinition.SetPlacementConstraints(input.PlacementConstraints)

	return taskDefinition
}

func (d *DeployProvider) createRenderedDeploy(body []byte) (*ecs.TaskDefinition, error) {
	var deployFile *ecs.TaskDefinition
	if err := json.Unmarshal(body, &deployFile); err != nil {
		return nil, fmt.Errorf("Failed to decode deploy: %s", err.Error())
	}

	if len(deployFile.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("Deploy must have at least one container definition")
	}

	for _, container := range deployFile.ContainerDefinitions {
		if container.LogConfiguration == nil {
			container.LogConfiguration = &ecs.LogConfiguration{
				LogDriver: aws.String("awslogs"),
				Options: map[string]*string{
					"awslogs-group":         aws.String(fmt.Sprintf("l0-%s", d.Config.Instance())),
					"awslogs-region":        aws.String(d.Config.Region()),
					"awslogs-stream-prefix": aws.String("l0"),
				},
			}
		}
	}

	return deployFile, nil
}
