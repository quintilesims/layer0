package provider

import "github.com/quintilesims/layer0/common/models"

type EnvironmentProvider interface {
	Create(req models.CreateEnvironmentRequest) (*models.Environment, error)
	Delete(environmentID string) error
	List() ([]models.EnvironmentSummary, error)
	Read(environmentID string) (*models.Environment, error)
	Update(req models.UpdateEnvironmentRequest) error
}
