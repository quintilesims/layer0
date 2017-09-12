package aws

import (
	"fmt"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/models"
)

func NewJobRunner(jobStore job.Store) job.RunnerFunc {
	return job.RunnerFunc(func(j models.Job) error {
		switch job.JobType(j.Type) {
		case job.DeleteEnvironmentJob:
			return deleteEnvironmentRunner(jobStore, j)
		default:
			return fmt.Errorf("Unrecognized JobType '%s'", j.Type)
		}

		return nil
	})
}

func deleteEnvironmentRunner(store job.Store, j models.Job) error {
	return nil
}
