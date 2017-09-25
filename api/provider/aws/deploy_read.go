package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (d *DeployProvider) Read(deployID string) (*models.Deploy, error) {
	fqDeployID := addLayer0Prefix(d.Config.Instance(), deployID)
	familyName := fqDeployID

	taskDefinitionOutput, err := d.describeTaskDefinition(familyName)
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

func (d *DeployProvider) describeTaskDefinition(familyName string) (*ecs.TaskDefinition, error) {
	input := &ecs.DescribeTaskDefinitionInput{}
	input.SetTaskDefinition(familyName)

	output, err := d.AWS.ECS.DescribeTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return output.TaskDefinition, nil
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
