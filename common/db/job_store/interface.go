package job_store

import (
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/models"
)

type JobStore interface {
	Init() error
	Delete(string) error
	Insert(*models.Job) error
	SelectAll() ([]*models.Job, error)
	SelectByID(string) (*models.Job, error)
	UpdateJobStatus(string, job.JobStatus) error
	SetJobMeta(string, map[string]string) error
}
