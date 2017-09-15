package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (d *DeployProvider) Create(req models.CreateDeployRequest) (*models.Deploy, error) {
	if req.DeployName == "" {
		return nil, fmt.Errorf("DeployName is required")
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

	return deploy, nil
}

func (d *DeployProvider) CreateDeploy(deployName string, body []byte) (*models.Deploy, error) {
	// since we use '.' as our ID-Version delimiter, we don't allow it in deploy names
	if strings.Contains(deployName, ".") {
		return nil, fmt.Errorf("Deploy names cannot contain '.'")
	}

	dockerrun, err := CreateRenderedDockerrun(body)
	if err != nil {
		return nil, err
	}

	// TODO: not convinced we need helper functions like these
	familyName := L0DeployID(deployName)
	if dockerrun.Family != "" && dockerrun.Family != familyName {
		return nil, fmt.Errorf("Custom family names are currently unsupported in Layer0")
	}

	taskDef := &ecs.RegisterTaskDefinitionInput{}
	taskDef.SetFamily(familyName)
	taskDef.SetTaskRoleArn(dockerrun.TaskRoleARN)
	taskDef.SetNetworkMode(dockerrun.NetworkMode)
	taskDef.SetContainerDefinitions(dockerrun.ContainerDefinitions)
	taskDef.SetVolumes(dockerrun.Volumes)
	// TODO: Understand why our model had to be updated to use ecs.TaskDefinitionPlacementConstraint
	// rather than ecs.PlacementConstraint
	taskDef.SetPlacementConstraints(dockerrun.PlacementConstraints)

	return d.populateModel(taskDef)
}

func (d *DeployProvider) populateModel(taskDef *ecs.RegisterTaskDefinitionInput) (*models.Deploy, error) {
	ecsDeployID := TaskDefinitionARNToECSDeployID(*taskDef.TaskRoleArn)

	containers := make([]*ecs.ContainerDefinition, len(taskDef.ContainerDefinitions))
	for i, c := range taskDef.ContainerDefinitions {
		containers[i] = &ecs.ContainerDefinition{c}
	}

	volumes := make([]*ecs.Volume, len(taskDef.Volumes))
	for i, v := range taskDef.Volumes {
		volumes[i] = &ecs.Volume{v}
	}

	placementConstraints := make([]*ecs.PlacementConstraint, len(taskDef.PlacementConstraints))
	for i, p := range taskDef.PlacementConstraints {
		placementConstraints[i] = &ecs.PlacementConstraint{p}
	}

	tmp := models.Dockerrun{
		ContainerDefinitions: containers,
		Volumes:              volumes,
		Family:               pstring(taskDef.Family),
		NetworkMode:          pstring(taskDef.NetworkMode),
		TaskRoleARN:          pstring(taskDef.TaskRoleArn),
		PlacementConstraints: placementConstraints,
	}

	dockerrun, err := json.Marshal(tmp)
	if err != nil {
		err := fmt.Errorf("Failed to extract dockerrun: %s", err.Error())
		return nil, errors.New(errors.InvalidJSON, err)
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
