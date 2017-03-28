package logic

import (
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type EnvironmentLogic interface {
	ListEnvironments() ([]*models.EnvironmentSummary, error)
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

func (e *L0EnvironmentLogic) ListEnvironments() ([]*models.EnvironmentSummary, error) {
	environments, err := e.Backend.ListEnvironments()
	if err != nil {
		return nil, err
	}

	summaries := make([]*models.EnvironmentSummary, len(environments))
	for i, environment := range environments {
		if err := e.populateModel(environment); err != nil {
			return nil, err
		}

		summaries[i] = &models.EnvironmentSummary{
			EnvironmentID:   environment.EnvironmentID,
			EnvironmentName: environment.EnvironmentName,
			OperatingSystem: environment.OperatingSystem,
		}
	}

	return summaries, nil
}

func (e *L0EnvironmentLogic) GetEnvironment(environmentID string) (*models.Environment, error) {
	environment, err := e.Backend.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}

	if err := e.populateModel(environment); err != nil {
		return nil, err
	}

	return environment, nil
}

func (e *L0EnvironmentLogic) DeleteEnvironment(environmentID string) error {
	if err := e.Backend.DeleteEnvironment(environmentID); err != nil {
		return err
	}

	if err := e.deleteEntityTags("environment", environmentID); err != nil {
		return err
	}

	return nil
}

func (e *L0EnvironmentLogic) CanCreateEnvironment(req models.CreateEnvironmentRequest) (bool, error) {
	tags, err := e.TagStore.SelectByQuery("environment", "")
	if err != nil {
		return false, err
	}

	tags = tags.WithKey("name").WithValue(req.EnvironmentName)
	return len(tags) == 0, nil
}

func (e *L0EnvironmentLogic) CreateEnvironment(req models.CreateEnvironmentRequest) (*models.Environment, error) {
	if req.EnvironmentName == "" {
		return nil, errors.Newf(errors.MissingParameter, "EnvironmentName is required")
	}

	if req.OperatingSystem == "" {
		return nil, errors.Newf(errors.MissingParameter, "OperatingSystem is required")
	}

	environment, err := e.Backend.CreateEnvironment(
		req.EnvironmentName,
		req.InstanceSize,
		req.OperatingSystem,
		req.MinClusterCount,
		req.UserDataTemplate)
	if err != nil {
		return nil, err
	}

	if err := e.upsertTagf(environment.EnvironmentID, "environment", "name", req.EnvironmentName); err != nil {
		return nil, err
	}

	if err := e.upsertTagf(environment.EnvironmentID, "environment", "os", req.OperatingSystem); err != nil {
		return nil, err
	}

	if err := e.populateModel(environment); err != nil {
		return environment, err
	}

	return environment, nil
}

func (e *L0EnvironmentLogic) UpdateEnvironment(environmentID string, minClusterCount int) (*models.Environment, error) {
	environment, err := e.Backend.UpdateEnvironment(environmentID, minClusterCount)
	if err != nil {
		return nil, err
	}

	if err := e.populateModel(environment); err != nil {
		return nil, err
	}

	return environment, nil
}

func (e *L0EnvironmentLogic) populateModel(model *models.Environment) error {
	tags, err := e.TagStore.SelectByQuery("environment", model.EnvironmentID)
	if err != nil {
		return err
	}

	if tag := tags.WithKey("name").First(); tag != nil {
		model.EnvironmentName = tag.Value
	}

	if tag := tags.WithKey("os").First(); tag != nil {
		model.OperatingSystem = tag.Value
	}

	return nil
}
