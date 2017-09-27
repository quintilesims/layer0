package aws

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
)

func (d *DeployProvider) Delete(deployID string) error {
	fqDeployID := addLayer0Prefix(d.Config.Instance(), deployID)

	taskArn, err := d.lookupTaskDefinitionARN(fqDeployID)
	if err != nil {
		return err
	}

	if err := d.deleteDeploy(taskArn); err != nil {
		return err
	}

	if err := d.deleteDeployTags(deployID); err != nil {
		return err
	}

	return nil
}

func (d *DeployProvider) deleteDeploy(taskArn string) error {
	input := &ecs.DeregisterTaskDefinitionInput{}
	input.SetTaskDefinition(taskArn)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := d.AWS.ECS.DeregisterTaskDefinition(input); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "does not exist" {
			return errors.Newf(errors.DeployDoesNotExist, "Deploy does not exist")
		}
	}
	return nil
}

func (d *DeployProvider) deleteDeployTags(deployID string) error {
	tags, err := d.TagStore.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := d.TagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
			return err
		}
	}

	return nil
}
