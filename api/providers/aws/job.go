package aws

import (
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/db/job_store"
	"github.com/quintilesims/layer0/common/models"
)

type Job struct {
	*AWSEntity
	JobStore job_store.JobStore
}

func NewJob(aws *awsc.Client, store job_store.JobStore, jobID string) *Job {
	return &Job{
		AWSEntity: NewAWSEntity(aws, jobID),
		JobStore:  store,
	}
}

func (j *Job) Delete() error {
	// todo: how to delete from a job scheduler?
	return j.JobStore.Delete(j.id)
}

func (j *Job) Model() (*models.Job, error) {
	return j.JobStore.SelectByID(j.id)
}
