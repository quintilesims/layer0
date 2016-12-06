package logic

import (
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

type AdminLogic interface {
	GetSQLStatus() (*models.SQLVersion, error)
	UpdateSQL() error
	GetHealth() (string, error)
}

type L0AdminLogic struct {
	Logic
}

func NewL0AdminLogic(lgc Logic) *L0AdminLogic {
	return &L0AdminLogic{lgc}
}

func (this *L0AdminLogic) GetSQLStatus() (*models.SQLVersion, error) {
	return this.AdminData.GetStatus()
}

func (this *L0AdminLogic) GetHealth() (string, error) {
	return this.Backend.GetRightSizerHealth()
}

func (this *L0AdminLogic) UpdateSQL() error {
	if err := this.AdminData.UpdateSQL(); err != nil {
		return err
	}

	return this.createDefaultTags()
}

func (this *L0AdminLogic) createDefaultTags() error {
	// certificate
	if err := this.upsertTagf(config.API_CERTIFICATE_ID, "certificate", "name", config.API_CERTIFICATE_NAME); err != nil {
		return err
	}

	// environment
	if err := this.upsertTagf(config.API_ENVIRONMENT_ID, "environment", "name", config.API_ENVIRONMENT_NAME); err != nil {
		return err
	}

	// load_balancer
	if err := this.upsertTagf(config.API_LOAD_BALANCER_ID, "load_balancer", "name", config.API_LOAD_BALANCER_NAME); err != nil {
		return err
	}

	if err := this.upsertTagf(config.API_LOAD_BALANCER_ID, "load_balancer", "environment_id", config.API_ENVIRONMENT_ID); err != nil {
		return err
	}

	// service
	if err := this.upsertTagf(config.API_SERVICE_ID, "service", "name", config.API_SERVICE_NAME); err != nil {
		return err
	}

	if err := this.upsertTagf(config.API_SERVICE_ID, "service", "environment_id", config.API_ENVIRONMENT_ID); err != nil {
		return err
	}

	if err := this.upsertTagf(config.API_SERVICE_ID, "service", "load_balancer_id", config.API_LOAD_BALANCER_ID); err != nil {
		return err
	}

	return nil
}
