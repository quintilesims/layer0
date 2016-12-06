package data

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
)

type JobData interface {
	ListJobs() ([]*models.Job, error)
	GetJob(jobID string) (*models.Job, error)
	UpdateJobStatus(jobID string, status types.JobStatus) error
	SetMeta(jobID, key, val string) error
	GetMeta(jobID string) (map[string]string, error)
	DeleteJob(jobID string) error
	CreateJob(job *models.Job) error
}

type JobLogicLayer struct {
	DataStore JobDataStore
}

func NewJobLogicLayer(dataStore JobDataStore) *JobLogicLayer {
	return &JobLogicLayer{
		DataStore: dataStore,
	}
}

func (this *JobLogicLayer) ListJobs() ([]*models.Job, error) {
	jobs, err := this.DataStore.Select()
	if err != nil {
		return nil, err
	}

	ret := make([]*models.Job, len(jobs))
	for i := range jobs {
		ret[i] = &jobs[i]
	}

	return ret, nil
}

func (this *JobLogicLayer) GetJob(jobID string) (*models.Job, error) {
	return this.DataStore.SelectByID(jobID)
}

func (this *JobLogicLayer) UpdateJobStatus(jobID string, status types.JobStatus) error {
	return this.DataStore.UpdateStatus(jobID, status)
}

func (this *JobLogicLayer) DeleteJob(jobID string) error {
	return this.DataStore.Delete(jobID)
}

func (this *JobLogicLayer) CreateJob(job *models.Job) error {
	return this.DataStore.Insert(job)
}

func (this *JobLogicLayer) SetMeta(jobID, key, val string) error {
	return this.DataStore.SetMeta(jobID, key, val)
}

func (this *JobLogicLayer) GetMeta(jobID string) (map[string]string, error) {
	return this.DataStore.GetMeta(jobID)
}
