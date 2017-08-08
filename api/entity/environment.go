package entity

import "github.com/quintilesims/layer0/common/models"

type Environment interface {
	Create(req models.CreateEnvironmentRequest) error
	Read() (*models.Environment, error)
	Delete() error
}
