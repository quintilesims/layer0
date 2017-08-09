package entity

import "github.com/quintilesims/layer0/common/models"

type Job interface {
	Model() (*models.Job, error)
	Delete() error
}
