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
	CreateEnvironmentLink(sourceEnvironmentID, destEnvironmentID string) error
	DeleteEnvironmentLink(sourceEnvironmentID, destEnvironmentID string) error
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
	tags, err := e.TagStore.SelectByQuery("environment", environmentID)
	if err != nil {
		return err
	}

	for _, tag := range tags.WithKey("link") {
		if err := e.DeleteEnvironmentLink(environmentID, tag.Value); err != nil {
			return err
		}
	}

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
		req.AMIID,
		req.MinClusterCount,
		req.UserDataTemplate)
	if err != nil {
		return nil, err
	}

	if err := e.upsertTag(models.Tag{EntityID: environment.EnvironmentID, EntityType: "environment", Key: "name", Value: req.EnvironmentName}); err != nil {
		return nil, err
	}

	if err := e.upsertTag(models.Tag{EntityID: environment.EnvironmentID, EntityType: "environment", Key: "os", Value: req.OperatingSystem}); err != nil {
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

func (e *L0EnvironmentLogic) CreateEnvironmentLink(sourceEnvironmentID, destEnvironmentID string) error {
	if err := e.Backend.CreateEnvironmentLink(sourceEnvironmentID, destEnvironmentID); err != nil {
		return nil
	}

	if err := e.upsertTag(models.Tag{EntityID: sourceEnvironmentID, EntityType: "environment", Key: "link", Value: destEnvironmentID}); err != nil {
		return nil
	}

	if err := e.upsertTag(models.Tag{EntityID: destEnvironmentID, EntityType: "environment", Key: "link", Value: sourceEnvironmentID}); err != nil {
		return nil
	}

	return nil
}

func (e *L0EnvironmentLogic) DeleteEnvironmentLink(sourceEnvironmentID, destEnvironmentID string) error {
	if err := e.Backend.DeleteEnvironmentLink(sourceEnvironmentID, destEnvironmentID); err != nil {
		return nil
	}

	sourceTags, err := e.TagStore.SelectByQuery("environment", sourceEnvironmentID)
	if err != nil {
		return err
	}

	for _, tag := range sourceTags.WithKey("link").WithValue(destEnvironmentID) {
		if err := e.TagStore.Delete(tag.TagID); err != nil {
			return err
		}
	}

	destTags, err := e.TagStore.SelectByQuery("environment", destEnvironmentID)
	if err != nil {
		return err
	}

	for _, tag := range destTags.WithKey("link").WithValue(sourceEnvironmentID) {
		if err := e.TagStore.Delete(tag.TagID); err != nil {
			return err
		}
	}

	return nil
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

	model.Links = []string{}
	for _, tag := range tags.WithKey("link") {
		model.Links = append(model.Links, tag.Value)
	}

	return nil
}
