package aws

import (
	"encoding/json"
	"fmt"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type JobRunner struct {
	deploy       provider.DeployProvider
	environment  provider.EnvironmentProvider
	loadBalancer provider.LoadBalancerProvider
	service      provider.ServiceProvider
	task         provider.TaskProvider
	jobStore     job.Store
}

func NewJobRunner(
	d provider.DeployProvider,
	e provider.EnvironmentProvider,
	l provider.LoadBalancerProvider,
	s provider.ServiceProvider,
	t provider.TaskProvider,
	store job.Store,
) *JobRunner {
	return &JobRunner{
		deploy:       d,
		environment:  e,
		loadBalancer: l,
		service:      s,
		task:         t,
		jobStore:     store,
	}
}

func (r *JobRunner) Run(j models.Job) error {
	switch job.JobType(j.Type) {
	case job.CreateEnvironmentJob:
		return r.createEnvironment(j)
	case job.DeleteEnvironmentJob:
		return r.deleteEnvironment(j)
	default:
		return fmt.Errorf("Unrecognized JobType '%s'", j.Type)
	}
}

/* Things to consider:
* When to retry
* When to fail
* Timeout
* Dependencies
 */

func (r *JobRunner) createEnvironment(j models.Job) error {
	var req models.CreateEnvironmentRequest
	if err := json.Unmarshal([]byte(j.Request), &req); err != nil {
		return errors.New(errors.InvalidRequest, err)
	}

	environment, err := r.environment.Create(req)
	if err != nil {
		return err
	}

	return r.jobStore.SetJobResult(j.JobID, environment.EnvironmentID)
}

func (r *JobRunner) deleteEnvironment(j models.Job) error {
	// todo: deps
	return r.environment.Delete(j.Request)
}
