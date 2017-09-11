package job

import (
	"github.com/quintilesims/layer0/common/models"
)

type Store interface {
	Init() error
	Delete(string) error
	Insert(models.Job) error
	SelectAll() ([]*models.Job, error)
	SelectByID(string) (*models.Job, error)
	UpdateStatus(string, Status) error
	SetJobMeta(string, map[string]string) error
}
