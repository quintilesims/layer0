package logic

import (
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

type EnvironmentLogic interface {
	ListEnvironments() ([]*models.Environment, error)
	GetEnvironment(id string) (*models.Environment, error)
	DeleteEnvironment(id string) error
	CanCreateEnvironment(req models.CreateEnvironmentRequest) (bool, error)
	CreateEnvironment(req models.CreateEnvironmentRequest) (*models.Environment, error)
	UpdateEnvironment(id string, minClusterCount int) (*models.Environment, error)
}

type L0EnvironmentLogic struct {
	Logic
}

func NewL0EnvironmentLogic(logic Logic) *L0EnvironmentLogic {
	return &L0EnvironmentLogic{
		Logic: logic,
	}
}

func (this *L0EnvironmentLogic) ListEnvironments() ([]*models.Environment, error) {
	environments, err := this.Backend.ListEnvironments()
	if err != nil {
		return nil, err
	}

	for _, environment := range environments {
		if err := this.populateModel(environment); err != nil {
			return nil, err
		}
	}

	return environments, nil
}

func (this *L0EnvironmentLogic) GetEnvironment(environmentID string) (*models.Environment, error) {
	environment, err := this.Backend.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}

	if err := this.populateModel(environment); err != nil {
		return nil, err
	}

	return environment, nil
}

func (this *L0EnvironmentLogic) DeleteEnvironment(environmentID string) error {
	if err := this.Backend.DeleteEnvironment(environmentID); err != nil {
		return err
	}

	if err := this.deleteEntityTags(environmentID, "environment"); err != nil {
		return err
	}

	return nil
}

func (this *L0EnvironmentLogic) CanCreateEnvironment(req models.CreateEnvironmentRequest) (bool, error) {
	filter := map[string]string{
		"type": "environment",
		"name": req.EnvironmentName,
	}

	tags, err := this.TagData.GetTags(filter)
	if err != nil {
		return false, err
	}

	return len(tags) == 0, nil
}

func (this *L0EnvironmentLogic) CreateEnvironment(req models.CreateEnvironmentRequest) (*models.Environment, error) {
	if req.EnvironmentName == "" {
		return nil, errors.Newf(errors.MissingParameter, "EnvironmentName is required")
	}

	environment, err := this.Backend.CreateEnvironment(req.EnvironmentName, req.InstanceSize, req.MinClusterCount, req.UserDataTemplate)
	if err != nil {
		return nil, err
	}

	if err := this.upsertTagf(environment.EnvironmentID, "environment", "name", req.EnvironmentName); err != nil {
		return nil, err
	}

	if err := this.populateModel(environment); err != nil {
		return environment, err
	}

	return environment, nil
}

func (this *L0EnvironmentLogic) UpdateEnvironment(environmentID string, minClusterCount int) (*models.Environment, error) {
	environment, err := this.Backend.UpdateEnvironment(environmentID, minClusterCount)
	if err != nil {
		return nil, err
	}

	if err := this.populateModel(environment); err != nil {
		return nil, err
	}

	return environment, nil
}

func (this *L0EnvironmentLogic) populateModel(model *models.Environment) error {
	filter := map[string]string{
		"type": "environment",
		"id":   model.EnvironmentID,
	}

	tags, err := this.TagData.GetTags(filter)
	if err != nil {
		return err
	}

	for _, tag := range rangeTags(tags) {
		if tag.Key == "name" {
			model.EnvironmentName = tag.Value
			break
		}
	}

	return nil
}
