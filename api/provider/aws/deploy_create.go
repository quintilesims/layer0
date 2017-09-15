package aws

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (d *DeployProvider) Create(req models.CreateDeployRequest) (*models.Deploy, error) {
	if req.DeployName == "" {
		return nil, fmt.Errorf("DeployName is required")
	}

	deploy, err := d.createDeploy(req.DeployName, req.Dockerrun)
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

func (d *DeployProvider) createDeploy(deployName string, body []byte) (*models.Deploy, error) {
	// since we use '.' as our ID-Version delimiter, we don't allow it in deploy names
	if strings.Contains(deployName, ".") {
		return nil, fmt.Errorf("Deploy names cannot contain '.'")
	}

	dockerrun, err := d.createRenderedDockerrun(body)
	if err != nil {
		return nil, err
	}

	familyName := L0DeployID(d.Config.Instance(), deployName)
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
		containers[i] = &ecs.ContainerDefinition{}
		containers[i] = c
	}

	volumes := make([]*ecs.Volume, len(taskDef.Volumes))
	for i, v := range taskDef.Volumes {
		volumes[i] = &ecs.Volume{}
		volumes[i] = v
	}

	placementConstraints := make([]*ecs.TaskDefinitionPlacementConstraint, len(taskDef.PlacementConstraints))
	for i, p := range taskDef.PlacementConstraints {
		placementConstraints[i] = &ecs.TaskDefinitionPlacementConstraint{}
		placementConstraints[i] = p
	}

	model := models.Dockerrun{
		ContainerDefinitions: containers,
		Volumes:              volumes,
		Family:               aws.StringValue(taskDef.Family),
		NetworkMode:          aws.StringValue(taskDef.NetworkMode),
		TaskRoleARN:          aws.StringValue(taskDef.TaskRoleArn),
		PlacementConstraints: placementConstraints,
	}

	dockerrun, err := json.Marshal(model)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract dockerrun: %s", err.Error())
	}

	deploy := &models.Deploy{
		DeployID:  L0DeployID(d.Config.Instance(), ecsDeployID),
		Version:   GetRevision(d.Config.Instance(), ecsDeployID),
		Dockerrun: dockerrun,
	}

	return deploy, nil
}

func (d *DeployProvider) createRenderedDockerrun(body []byte) (*models.Dockerrun, error) {
	var dockerrun *models.Dockerrun
	if err := json.Unmarshal(body, &dockerrun); err != nil {
		return nil, fmt.Errorf("Failed to decode deploy: %s", err.Error())
	}

	if len(dockerrun.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("Deploy must have at least one container definition")
	}

	for _, container := range dockerrun.ContainerDefinitions {
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

	return dockerrun, nil
}
