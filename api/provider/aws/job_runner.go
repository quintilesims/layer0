package aws

import (
	"fmt"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/models"
)

type JobRunner struct {
	store job.JobStore
}

func NewJobRunner(s job.JobStore) *JobRunner {
	return &JobRunner{
		store: s,
	}
}

func (r *JobRunner) Run(j models.Job) error {
	switch job.JobType(j.JobType) {
	case job.DeleteEnvironmentJob:
		// todo: run + retry
	default:
		return fmt.Errorf("Unrecognized JobType '%s'", j.JobType)
	}

	return nil
}
