package logic

import (
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

type AdminLogic interface {
	RunEnvironmentScaler(string) (*models.ScalerRunInfo, error)
	UpdateSQL() error
}

type L0AdminLogic struct {
	Logic
}

func NewL0AdminLogic(l Logic) *L0AdminLogic {
	return &L0AdminLogic{
		Logic: l,
	}
}

func (a *L0AdminLogic) RunEnvironmentScaler(environmentID string) (*models.ScalerRunInfo, error) {
	return a.Logic.Scaler.Scale(environmentID)
}

func (a *L0AdminLogic) UpdateSQL() error {
	if err := a.TagStore.Init(); err != nil {
		return err
	}

	if err := a.JobStore.Init(); err != nil {
		return err
	}

	return a.createDefaultTags()
}

func (a *L0AdminLogic) createDefaultTags() error {
	// environment
	if err := a.TagStore.Insert(models.Tag{EntityID: config.API_ENVIRONMENT_ID, EntityType: "environment", Key: "name", Value: config.API_ENVIRONMENT_NAME}); err != nil {
		return err
	}

	if err := a.TagStore.Insert(models.Tag{EntityID: config.API_ENVIRONMENT_ID, EntityType: "environment", Key: "os", Value: "linux"}); err != nil {
		return err
	}

	// load_balancer
	if err := a.TagStore.Insert(models.Tag{EntityID: config.API_LOAD_BALANCER_ID, EntityType: "load_balancer", Key: "name", Value: config.API_LOAD_BALANCER_NAME}); err != nil {
		return err
	}

	if err := a.TagStore.Insert(models.Tag{EntityID: config.API_LOAD_BALANCER_ID, EntityType: "load_balancer", Key: "environment_id", Value: config.API_ENVIRONMENT_ID}); err != nil {
		return err
	}

	// service
	if err := a.TagStore.Insert(models.Tag{EntityID: config.API_SERVICE_ID, EntityType: "service", Key: "name", Value: config.API_SERVICE_NAME}); err != nil {
		return err
	}

	if err := a.TagStore.Insert(models.Tag{EntityID: config.API_SERVICE_ID, EntityType: "service", Key: "environment_id", Value: config.API_ENVIRONMENT_ID}); err != nil {
		return err
	}

	if err := a.TagStore.Insert(models.Tag{EntityID: config.API_SERVICE_ID, EntityType: "service", Key: "load_balancer_id", Value: config.API_LOAD_BALANCER_ID}); err != nil {
		return err
	}

	return nil
}
