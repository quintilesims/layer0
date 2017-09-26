package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func (d *DeployProvider) Read(deployID string) (*models.Deploy, error) {
	taskDefinitionARN, err := d.lookupTaskDefinitionARN(deployID)
	if err != nil {
		return nil, err
	}

	taskDefinitionOutput, err := d.describeTaskDefinition(taskDefinitionARN)
	if err != nil {
		return nil, err
	}

	deployFile, err := json.Marshal(taskDefinitionOutput)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract deploy file: %s", err.Error())
	}

	model := &models.Deploy{
		DeployID:   deployID,
		DeployFile: deployFile,
	}

	if err := d.populateModelTags(deployID, model); err != nil {
		return nil, err
	}

	return model, nil
}

func (d *DeployProvider) describeTaskDefinition(taskDefinitionARN string) (*ecs.TaskDefinition, error) {
	input := &ecs.DescribeTaskDefinitionInput{}
	input.SetTaskDefinition(taskDefinitionARN)

	output, err := d.AWS.ECS.DescribeTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return output.TaskDefinition, nil
}

func (d *DeployProvider) lookupTaskDefinitionARN(deployID string) (string, error) {
	tags, err := d.TagStore.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", errors.Newf(errors.DeployDoesNotExist, "Deploy '%s' does not exist", deployID)
	}

	if tag, ok := tags.WithKey("arn").First(); ok {
		return tag.Value, nil
	}

	return "", fmt.Errorf("Failed to find ARN for deploy '%s'", deployID)
}

func (d *DeployProvider) populateModelTags(deployID string, model *models.Deploy) error {
	tags, err := d.TagStore.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return err
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.DeployName = tag.Value
	}

	if tag, ok := tags.WithKey("version").First(); ok {
		model.Version = tag.Value
	}

	return nil
}
