package provider

import "github.com/quintilesims/layer0/common/models"

type EnvironmentProvider interface {
	Create(req models.CreateEnvironmentRequest) (*models.Environment, error)
	Read(environmentID string) (*models.Environment, error)
	List() ([]models.EnvironmentSummary, error)
	Update(req models.UpdateEnvironmentRequest) error
	Delete(environmentID string) error
}
