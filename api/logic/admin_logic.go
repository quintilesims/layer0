package logic

import (
	"github.com/quintilesims/layer0/common/config"
)

type AdminLogic interface {
	UpdateSQL() error
	GetHealth() (string, error)
	RunRightSizer() error
}

type L0AdminLogic struct {
	Logic
}

func NewL0AdminLogic(lgc Logic) *L0AdminLogic {
	return &L0AdminLogic{lgc}
}

func (a *L0AdminLogic) GetHealth() (string, error) {
	return a.Backend.GetRightSizerHealth()
}

func (a L0AdminLogic) RunRightSizer() error {
	return	a.Backend.RunRightSizer()
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
