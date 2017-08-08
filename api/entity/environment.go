package entity

import "github.com/quintilesims/layer0/common/models"

type Environment interface {
	Create(req models.CreateEnvironmentRequest) error
	Model() (*models.Environment, error)
	Summary() (*models.EnvironmentSummary, error)
	Delete() error
}
