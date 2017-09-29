package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

// Read returns a models.Deploy based on the provided Deploy ID. The Deploy ID
// is used to look up the associated Task Definition ARN. The Task Definition
// ARN is used when the DescribeTaskDefinition request is made to AWS.
func (d *DeployProvider) Read(deployID string) (*models.Deploy, error) {
	deployModel, err := d.newDeployModel(deployID)
	if err != nil {
		return nil, err
	}

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

	deployModel.DeployFile = deployFile

	return deployModel, nil
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

func (d *DeployProvider) newDeployModel(deployID string) (*models.Deploy, error) {
	model := &models.Deploy{
		DeployID: deployID,
	}

	tags, err := d.TagStore.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return nil, err
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.DeployName = tag.Value
	}

	if tag, ok := tags.WithKey("version").First(); ok {
		model.Version = tag.Value
	}

	return model, nil
}
