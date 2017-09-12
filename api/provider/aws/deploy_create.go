package aws

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (d *DeployProvider) Create(req models.CreateDeployRequest) (*models.Deploy, error) {
	if req.DeployName == "" {
		return nil, errors.Newf(errors.MissingParameter, "DeployName is required")
	}

	deploy, err := d.CreateDeploy(req.DeployName, req.Dockerrun)
	if err != nil {
		return deploy, err
	}

	if err := d.TagStore.Insert(models.Tag{EntityID: deploy.DeployID, EntityType: "deploy", Key: "name", Value: req.DeployName}); err != nil {
		return deploy, err
	}

	if err := d.TagStore.Insert(models.Tag{EntityID: deploy.DeployID, EntityType: "deploy", Key: "version", Value: deploy.Version}); err != nil {
		return deploy, err
	}

	if err := d.populateModel(deploy); err != nil {
		return deploy, err
	}

	return deploy, nil
}

func (d *DeployProvider) CreateDeploy(deployName string, body []byte) (*models.Deploy, error) {
	// since we use '.' as our ID-Version delimiter, we don't allow it in deploy names
	if strings.Contains(deployName, ".") {
		return nil, errors.Newf(errors.InvalidDeployID, "Deploy names cannot contain '.'")
	}

	dockerrun, err := CreateRenderedDockerrun(body)
	if err != nil {
		return nil, err
	}

	// deploys for jobs will have the deploy.Family field, but they will match familyName
	familyName := id.L0DeployID(deployName).ECSDeployID().String()
	if dockerrun.Family != "" && dockerrun.Family != familyName {
		return nil, fmt.Errorf("Custom family names are currently unsupported in Layer0")
	}

	taskDef, err := ecs.RegisterTaskDefinitionInput(
		familyName,
		dockerrun.TaskRoleARN,
		dockerrun.NetworkMode,
		dockerrun.ContainerDefinitions,
		dockerrun.Volumes,
		dockerrun.PlacementConstraints)
	if err != nil {
		return nil, err
	}

	return this.populateModel(taskDef)
}

func (d *DeployProvider) populateModel(taskDef *ecs.TaskDefinition) (*models.Deploy, error) {
	ecsDeployID := id.TaskDefinitionARNToECSDeployID(*taskDef.TaskDefinitionArn)

	dockerrun, err := extractDockerrun(taskDef)
	if err != nil {
		return nil, err
	}

	deploy := &models.Deploy{
		DeployID:  ecsDeployID.L0DeployID(),
		Version:   ecsDeployID.Revision(),
		Dockerrun: dockerrun,
	}

	return deploy, nil
}

var CreateRenderedDockerrun = func(body []byte) (*models.Dockerrun, error) {
	dockerrun, err := MarshalDockerrun(body)
	if err != nil {
		return nil, err
	}

	for _, container := range dockerrun.ContainerDefinitions {
		if container.LogConfiguration == nil {
			container.LogConfiguration = &awsecs.LogConfiguration{
				LogDriver: stringp("awslogs"),
				Options: map[string]*string{
					"awslogs-group":         stringp(config.AWSLogGroupID()),
					"awslogs-region":        stringp(config.AWSRegion()),
					"awslogs-stream-prefix": stringp("l0"),
				},
			}
		}
	}

	return dockerrun, nil
}
