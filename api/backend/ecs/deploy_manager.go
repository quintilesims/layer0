package ecsbackend

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type ECSDeployManager struct {
	ECS ecs.Provider
}

func NewECSDeployManager(ecsprovider ecs.Provider) *ECSDeployManager {
	return &ECSDeployManager{
		ECS: ecsprovider,
	}
}

func (this *ECSDeployManager) ListDeploys() ([]*models.Deploy, error) {
	taskDefinitionARNs, err := this.ECS.Helper_ListTaskDefinitions(id.PREFIX)
	if err != nil {
		return nil, err
	}

	deploys := make([]*models.Deploy, len(taskDefinitionARNs))
	for i, taskDefinitionARN := range taskDefinitionARNs {
		ecsDeployID := id.TaskDefinitionARNToECSDeployID(*taskDefinitionARN)
		deploys[i] = &models.Deploy{
			DeployID: ecsDeployID.L0DeployID(),
		}
	}

	return deploys, nil
}

func (this *ECSDeployManager) GetDeploy(deployID string) (*models.Deploy, error) {
	ecsDeployID := id.L0DeployID(deployID).ECSDeployID()

	taskDef, err := this.ECS.DescribeTaskDefinition(ecsDeployID.TaskDefinition())
	if err != nil {
		if ContainsErrMsg(err, "Unable to describe task definition") {
			err := fmt.Errorf("Deploy with id '%s' does not exist", deployID)
			return nil, errors.New(errors.InvalidDeployID, err)
		}

		return nil, err
	}

	return this.populateModel(taskDef)
}

func (this *ECSDeployManager) DeleteDeploy(deployID string) error {
	ecsDeployID := id.L0DeployID(deployID).ECSDeployID()

	if err := this.ECS.DeleteTaskDefinition(ecsDeployID.TaskDefinition()); err != nil {
		if ContainsErrMsg(err, "does not exist") {
			err := fmt.Errorf("Deploy with id '%s' does not exist", deployID)
			return errors.New(errors.InvalidDeployID, err)
		}

		return err
	}

	return nil
}

func (this *ECSDeployManager) CreateDeploy(deployName string, body []byte) (*models.Deploy, error) {
	// since we use '.' as our ID-Version delimiter, we don't allow it in deploy names
	if strings.Contains(deployName, ".") {
		return nil, errors.Newf(errors.InvalidDeployID, "Deploy names cannot contain '.'")
	}

	deploy, err := CreateRenderedDeploy(body)
	if err != nil {
		return nil, err
	}

	// deploys for jobs will have the deploy.Family field, but they will match familyName
	familyName := id.L0DeployID(deployName).ECSDeployID().String()
	if deploy.Family != "" && deploy.Family != familyName {
		return nil, fmt.Errorf("Custom family names are currently unsupported in Layer0")
	}

	taskDef, err := this.ECS.RegisterTaskDefinition(
		familyName,
		deploy.TaskRoleARN,
		deploy.NetworkMode,
		deploy.ContainerDefinitions,
		deploy.Volumes,
		deploy.PlacementConstraints)
	if err != nil {
		return nil, err
	}

	return this.populateModel(taskDef)
}

func (this *ECSDeployManager) populateModel(taskDef *ecs.TaskDefinition) (*models.Deploy, error) {
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

type Deploy struct {
	ContainerDefinitions []*ecs.ContainerDefinition `json:"containerDefinitions,omitempty"`
	Volumes              []*ecs.Volume              `json:"volumes,omitempty"`
	Family               string                     `json:"family,omitempty"`
	NetworkMode          string                     `json:"networkMode,omitempty"`
	TaskRoleARN          string                     `json:"taskRoleArn,omitempty"`
	PlacementConstraints []*ecs.PlacementConstraint `json:"placementConstraints,omitempty"`
}

func MarshalDeploy(body []byte) (*Deploy, error) {
	var d Deploy
	if err := json.Unmarshal(body, &d); err != nil {
		err := fmt.Errorf("Failed to decode deploy: %s", err.Error())
		return nil, errors.New(errors.InvalidJSON, err)
	}

	if len(d.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("Deploy must have at least one container definition")
	}

	return &d, nil
}

func extractDockerrun(taskDef *ecs.TaskDefinition) ([]byte, error) {
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

	d := Deploy{
		ContainerDefinitions: containers,
		Volumes:              volumes,
		Family:               pstring(taskDef.Family),
		NetworkMode:          pstring(taskDef.NetworkMode),
		TaskRoleARN:          pstring(taskDef.TaskRoleArn),
		PlacementConstraints: placementConstraints,
	}

	bytes, err := json.Marshal(d)
	if err != nil {
		err := fmt.Errorf("Failed to extract dockerrun: %s", err.Error())
		return nil, errors.New(errors.InvalidJSON, err)
	}

	return bytes, nil
}
