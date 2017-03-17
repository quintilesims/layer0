package logic

import (
	"github.com/quintilesims/layer0/api/scheduler/resource"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

type AdminLogic interface {
	UpdateSQL() error
	RunResourceManager(environmentID string) (*models.ScalerRunInfo, error)
}

type L0AdminLogic struct {
	Logic
	ResourceManager *resource.ResourceManager
}

func NewL0AdminLogic(l Logic, r *resource.ResourceManager) *L0AdminLogic {
	return &L0AdminLogic{
		Logic:           l,
		ResourceManager: r,
	}
}

func (a L0AdminLogic) RunResourceManager(environmentID string) (*models.ScalerRunInfo, error) {
	return a.ResourceManager.Run(environmentID)
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
	if err := a.upsertTagf(config.API_ENVIRONMENT_ID, "environment", "name", config.API_ENVIRONMENT_NAME); err != nil {
		return err
	}

	// load_balancer
	if err := a.upsertTagf(config.API_LOAD_BALANCER_ID, "load_balancer", "name", config.API_LOAD_BALANCER_NAME); err != nil {
		return err
	}

	if err := a.upsertTagf(config.API_LOAD_BALANCER_ID, "load_balancer", "environment_id", config.API_ENVIRONMENT_ID); err != nil {
		return err
	}

	// service
	if err := a.upsertTagf(config.API_SERVICE_ID, "service", "name", config.API_SERVICE_NAME); err != nil {
		return err
	}

	if err := a.upsertTagf(config.API_SERVICE_ID, "service", "environment_id", config.API_ENVIRONMENT_ID); err != nil {
		return err
	}

	if err := a.upsertTagf(config.API_SERVICE_ID, "service", "load_balancer_id", config.API_LOAD_BALANCER_ID); err != nil {
		return err
	}

	return nil
}
