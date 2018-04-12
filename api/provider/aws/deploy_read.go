package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

// Read returns a *models.Deploy based on the provided deployID. The deployID
// is used to look up the associated Task Definition ARN. The Task Definition
// ARN is subsequently used when the DescribeTaskDefinition request is made to AWS.
func (d *DeployProvider) Read(deployID string) (*models.Deploy, error) {
	taskDefinitionARN, err := d.lookupTaskDefinitionARN(deployID)
	if err != nil {
		return nil, err
	}

	taskDefinition, err := describeTaskDefinition(d.AWS.ECS, taskDefinitionARN)
	if err != nil {
		return nil, err
	}

	deployCompatibilities := d.extractDeployCompatibilities(taskDefinition)

	deployFile, err := json.Marshal(taskDefinition)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract deploy file: %s", err.Error())
	}

	return d.makeDeployModel(deployID, deployCompatibilities, deployFile)
}

func (d *DeployProvider) extractDeployCompatibilities(taskDefinition *ecs.TaskDefinition) []string {
	compatibilities := []string{}
	for _, compatibility := range taskDefinition.Compatibilities {
		switch aws.StringValue(compatibility) {
		case ecs.LaunchTypeEc2:
			compatibilities = append(compatibilities, models.DeployCompatibilityStateful)
		case ecs.LaunchTypeFargate:
			compatibilities = append(compatibilities, models.DeployCompatibilityStateless)
		}
	}

	return compatibilities
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

func (d *DeployProvider) makeDeployModel(deployID string, deployCompatibilities []string, deployFile []byte) (*models.Deploy, error) {
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

	model.Compatibilities = deployCompatibilities
	model.DeployFile = deployFile

	return model, nil
}
