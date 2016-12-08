package data

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
)

// JobDataStore ...
type JobDataStore interface {
	Select() ([]models.Job, error)
	SelectByID(id string) (*models.Job, error)
	Insert(job *models.Job) error
	Delete(id string) error
	UpdateStatus(jobID string, jobStatus types.JobStatus) error
	SetMeta(jobID, key, val string) error
	GetMeta(jobID string) (map[string]string, error)
}
