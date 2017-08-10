package aws

import (
	"github.com/quintilesims/layer0/api/entity"
	"github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/db/job_store"
)

type AWSProvider struct {
	AWS      *aws.Client
	JobStore job_store.JobStore
	// todo: Tag tag.Provider
	// todo: scheduler
}

func NewAWSProvider(a *aws.Client, j job_store.JobStore) *AWSProvider {
	return &AWSProvider{
		AWS:      a,
		JobStore: j,
	}
}

func (a *AWSProvider) ListEnvironmentIDs() ([]string, error) {
	return nil, nil
}

func (a *AWSProvider) GetEnvironment(environmentID string) entity.Environment {
	return NewEnvironment(a.AWS, environmentID)
}

func (a *AWSProvider) ListJobIDs() ([]string, error) {
	return nil, nil
}

func (a *AWSProvider) GetJob(jobID string) entity.Job {
	return NewJob(a.AWS, a.JobStore, jobID)
}
