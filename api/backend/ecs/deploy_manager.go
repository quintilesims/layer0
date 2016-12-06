package ecsbackend

import (
	"encoding/json"
	"fmt"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs/id"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs"
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"strings"
)

type ECSDeployManager struct {
	ECS           ecs.Provider
	ClusterScaler ClusterScaler
}

func NewECSDeployManager(ecsprovider ecs.Provider, cluster ClusterScaler) *ECSDeployManager {
	return &ECSDeployManager{
		ecsprovider,
		cluster,
	}
}

func (this *ECSDeployManager) ListDeploys() ([]*models.Deploy, error) {
	taskDefinitionARNs, err := this.ECS.Helper_ListTaskDefinitions(id.PREFIX)
	if err != nil {
		return nil, err
	}

	models := []*models.Deploy{}
	for _, taskDefinitionARN := range taskDefinitionARNs {
		ecsDeployID := id.TaskDefinitionARNToECSDeployID(*taskDefinitionARN)

		if strings.HasPrefix(ecsDeployID.String(), id.PREFIX) {
			model := this.populateModel(ecsDeployID)
			models = append(models, model)
		}
	}

	return models, nil
}

func (this *ECSDeployManager) GetDeploy(deployID string) (*models.Deploy, error) {
	ecsDeployID := id.L0DeployID(deployID).ECSDeployID()

	if _, err := this.ECS.DescribeTaskDefinition(ecsDeployID.TaskDefinition()); err != nil {
		if ContainsErrMsg(err, "Unable to describe task definition") {
			err := fmt.Errorf("Deploy with id '%s' does not exist", deployID)
			return nil, errors.New(errors.InvalidDeployID, err)
		}

		return nil, err
	}

	return this.populateModel(ecsDeployID), nil
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

	taskDefinition, err := this.ECS.RegisterTaskDefinition(
		familyName,
		deploy.TaskRoleARN,
		deploy.NetworkMode,
		deploy.ContainerDefinitions,
		deploy.Volumes)
	if err != nil {
		return nil, err
	}

	// unlike other entities, let aws generate our entity id
	// ecs returns <deployName>:<revision> as the unique deploy id
	ecsDeployID := id.TaskDefinitionARNToECSDeployID(*taskDefinition.TaskDefinitionArn)
	return this.populateModel(ecsDeployID), nil
}

func (this *ECSDeployManager) populateModel(ecsDeployID id.ECSDeployID) *models.Deploy {
	return &models.Deploy{
		DeployID: ecsDeployID.L0DeployID(),
		Version:  ecsDeployID.Revision(),
	}
}

type deploy struct {
	ContainerDefinitions []*ecs.ContainerDefinition `json:"containerDefinitions,omitempty"`
	Volumes              []*ecs.Volume              `json:"volumes,omitempty"`
	Family               string                     `json:"family,omitempty"`
	NetworkMode          string                     `json:"networkMode,omitempty"`
	TaskRoleARN          string                     `json:"taskRoleArn,omitempty"`
}

func marshalDeploy(body []byte) (*deploy, error) {
	var d deploy
	if err := json.Unmarshal(body, &d); err != nil {
		err := fmt.Errorf("Failed to decode deploy: %s", err.Error())
		return nil, errors.New(errors.InvalidJSON, err)
	}

	if len(d.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("Deploy must have at least one container definition")
	}

	return &d, nil
}
