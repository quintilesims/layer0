package aws

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// Create takes a request from the Controller layer, formats the request,
// and calls the AWS ECS API with it. In order to do this, we must:
//
// - Generate a unique deployID and fqDeployID
// - Render a AWS TaskDefinition from the Request
// - Create an ecs.TaskDefinition object for the Task Definition data
// - Call the AWS' ECS RegisterTaskDefinition API with the previous steps' object
// - If the AWS API call is successful, Marshal the request into a []byte
//   - Instantiate Layer0 models.Deploy object with the marshaled data
// - Insert entity tags into Layer0 tags database
// - Return entity object back to Controller layer
//
func (d *DeployProvider) Create(req models.CreateDeployRequest) (*models.Deploy, error) {
	deployID := generateEntityID(req.DeployName)
	fqDeployID := addLayer0Prefix(d.Config.Instance(), deployID)
	familyName := fqDeployID

	taskDefinition, err := d.renderTaskDefinition(req.DeployFile, fqDeployID)
	if err != nil {
		return nil, err
	}

	input := &ecs.TaskDefinition{}
	input.SetFamily(familyName)
	input.SetTaskRoleArn(aws.StringValue(taskDefinition.TaskRoleArn))
	input.SetNetworkMode(aws.StringValue(taskDefinition.NetworkMode))
	input.SetContainerDefinitions(taskDefinition.ContainerDefinitions)
	input.SetVolumes(taskDefinition.Volumes)
	input.SetPlacementConstraints(taskDefinition.PlacementConstraints)

	taskDefinitionOutput, err := d.createTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(taskDefinitionOutput)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract deploy file: %s", err.Error())
	}

	deploy := &models.Deploy{
		// TODO: we previously removed prefix from deployID ala below. Is this
		// no longer needed?
		// func removePrefix(id string) string {
		// 	return strings.TrimPrefix(id, PREFIX)
		// }
		DeployID:   deployID,
		Version:    d.getDeployRevision(deployID),
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

	_, err := d.AWS.ECS.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	// TODO: return original request or RegisterTaskDefinition() input?
	return taskDefinitionRequest, nil
}

func (d *DeployProvider) renderTaskDefinition(body []byte, familyName string) (*ecs.TaskDefinition, error) {
	var taskDefinition *ecs.TaskDefinition

	if err := json.Unmarshal(body, &taskDefinition); err != nil {
		return nil, fmt.Errorf("Failed to decode deploy: %s", err.Error())
	}

	if len(taskDefinition.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("Deploy must have at least one container definition")
	}

	if taskDefinition.Family != nil && taskDefinition.Family != aws.String(familyName) {
		return nil, fmt.Errorf("Custom family names are currently unsupported in Layer0")
	}

	for _, container := range taskDefinition.ContainerDefinitions {
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

	return taskDefinition, nil
}

// TODO: remove this if possible
func (d *DeployProvider) getDeployRevision(id string) string {
	if split := strings.Split(id, "."); len(split) > 1 {
		return split[1]
	}

	return "1"
}
