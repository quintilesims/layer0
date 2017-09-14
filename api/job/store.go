package job

import (
	"github.com/quintilesims/layer0/common/models"
)

type Store interface {
	Delete(jobID string) error
	AcquireJob(jobID string) (bool, error)
	Insert(jobType JobType, req string) (string, error)
	SetInsertHook(hook func(jobID string))
	SelectAll() ([]*models.Job, error)
	SelectByID(jobID string) (*models.Job, error)
	SetJobStatus(jobID string, status Status) error
	SetJobResult(jobID, result string) error
	SetJobError(jobID string, err error) error
}
