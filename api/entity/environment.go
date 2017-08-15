package entity

import "github.com/quintilesims/layer0/common/models"

type EnvironmentProvider interface {
	Create(req models.CreateEnvironmentRequest) error
	Read(environmentID string) (*models.Environment, error)
	List() ([]*models.EnvironmentSummary, error)
	Delete(environmentID string) error
}
