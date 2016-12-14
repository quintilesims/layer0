package logic

import (
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
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

	if err := this.deleteEntityTags("environment", environmentID); err != nil {
		return err
	}

	return nil
}

func (this *L0EnvironmentLogic) CanCreateEnvironment(req models.CreateEnvironmentRequest) (bool, error) {
	tags, err := this.TagStore.SelectByQuery("environment", "")
	if err != nil {
		return false, err
	}

	tags = tags.WithKey("name").WithValue(req.EnvironmentName)
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
	tags, err := this.TagStore.SelectByQuery("environment", model.EnvironmentID)
	if err != nil {
		return err
	}

	if tag := tags.WithKey("name").First(); tag != nil {
		model.EnvironmentName = tag.Value
	}

	return nil
}
