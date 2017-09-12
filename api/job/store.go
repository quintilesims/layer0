package job

import (
	"github.com/quintilesims/layer0/common/models"
)

type Store interface {
	Delete(jobID string) error
	Insert(jobType JobType, req string) (string, error)
	SelectAll() ([]*models.Job, error)
	SelectByID(jobID string) (*models.Job, error)
	SetJobStatus(jobID string, status Status) error
	SetJobMeta(jobID string, meta map[string]string) error
	SetJobError(jobID string, err error) error
}
