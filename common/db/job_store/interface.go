package job_store

import (
	"github.com/quintilesims/layer0/common/models"
)

type JobStore interface {
	Init() error
	Close()
	Delete(string) error
	Insert(*models.Job) error
	SelectAll() ([]*models.Job, error)
	SelectByID(string) (*models.Job, error)
}
