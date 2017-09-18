package aws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/api/scaler"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type JobRunner struct {
	deployProvider       provider.DeployProvider
	environmentProvider  provider.EnvironmentProvider
	loadBalancerProvider provider.LoadBalancerProvider
	serviceProvider      provider.ServiceProvider
	taskProvider         provider.TaskProvider
	scaler               *scaler.Dispatcher
	jobStore             job.Store
}

func NewJobRunner(
	d provider.DeployProvider,
	e provider.EnvironmentProvider,
	l provider.LoadBalancerProvider,
	s provider.ServiceProvider,
	t provider.TaskProvider,
	scaler *scaler.Dispatcher,
	store job.Store,
) *JobRunner {
	return &JobRunner{
		deployProvider:       d,
		environmentProvider:  e,
		loadBalancerProvider: l,
		serviceProvider:      s,
		taskProvider:         t,
		scaler:               scaler,
		jobStore:             store,
	}
}

func (r *JobRunner) Run(j models.Job) error {
	switch job.JobType(j.Type) {
	case job.CreateDeployJob:
		return r.createDeploy(j.JobID, j.Request)
	case job.CreateEnvironmentJob:
		return r.createEnvironment(j.JobID, j.Request)
	case job.CreateLoadBalancerJob:
		return r.createLoadBalancer(j.JobID, j.Request)
	case job.CreateServiceJob:
		return r.createService(j.JobID, j.Request)
	case job.CreateTaskJob:
		return r.createTask(j.JobID, j.Request)
	case job.DeleteDeployJob:
		return r.deleteDeploy(j.JobID, j.Request)
	case job.DeleteEnvironmentJob:
		return r.deleteEnvironment(j.JobID, j.Request)
	case job.DeleteLoadBalancerJob:
		return r.deleteLoadBalancer(j.JobID, j.Request)
	case job.DeleteServiceJob:
		return r.deleteService(j.JobID, j.Request)
	case job.DeleteTaskJob:
		return r.deleteTask(j.JobID, j.Request)
	case job.UpdateEnvironmentJob:
		return r.updateEnvironment(j.JobID, j.Request)
	case job.UpdateLoadBalancerJob:
		return r.updateLoadBalancer(j.JobID, j.Request)
	case job.UpdateServiceJob:
		return r.updateService(j.JobID, j.Request)
	default:
		return fmt.Errorf("Unrecognized JobType '%s'", j.Type)
	}
}

func (r *JobRunner) createDeploy(jobID, request string) error {
	var req models.CreateDeployRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return errors.New(errors.InvalidRequest, err)
	}

	deploy, err := r.deployProvider.Create(req)
	if err != nil {
		return err
	}

	return r.jobStore.SetJobResult(jobID, deploy.DeployID)
}

func (r *JobRunner) createEnvironment(jobID, request string) error {
	var req models.CreateEnvironmentRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return errors.New(errors.InvalidRequest, err)
	}

	environment, err := r.environmentProvider.Create(req)
	if err != nil {
		return err
	}

	return r.jobStore.SetJobResult(jobID, environment.EnvironmentID)
}

func (r *JobRunner) createLoadBalancer(jobID, request string) error {
	var req models.CreateLoadBalancerRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return errors.New(errors.InvalidRequest, err)
	}

	loadBalancer, err := r.loadBalancerProvider.Create(req)
	if err != nil {
		return err
	}

	return r.jobStore.SetJobResult(jobID, loadBalancer.LoadBalancerID)
}

func (r *JobRunner) createService(jobID, request string) error {
	var req models.CreateServiceRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return errors.New(errors.InvalidRequest, err)
	}

	r.scaler.ScheduleRun(req.EnvironmentID)

	service, err := r.serviceProvider.Create(req)
	if err != nil {
		return err
	}

	return r.jobStore.SetJobResult(jobID, service.ServiceID)
}

func (r *JobRunner) createTask(jobID, request string) error {
	var req models.CreateTaskRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return errors.New(errors.InvalidRequest, err)
	}

	r.scaler.ScheduleRun(req.EnvironmentID)

	task, err := r.taskProvider.Create(req)
	if err != nil {
		return err
	}

	return r.jobStore.SetJobResult(jobID, task.TaskID)
}

func (r *JobRunner) deleteDeploy(jobID, deployID string) error {
	return r.deployProvider.Delete(deployID)
}

func (r *JobRunner) deleteEnvironment(jobID, environmentID string) error {
	loadBalancers, err := r.loadBalancerProvider.List()
	if err != nil {
		return err
	}

	for _, loadBalancer := range loadBalancers {
		if loadBalancer.EnvironmentID == environmentID {
			if err := r.deleteLoadBalancer(jobID, loadBalancer.LoadBalancerID); err != nil {
				return err
			}
		}
	}

	services, err := r.serviceProvider.List()
	if err != nil {
		return err
	}

	for _, service := range services {
		if service.EnvironmentID == environmentID {
			if err := r.deleteService(jobID, service.ServiceID); err != nil {
				return err
			}
		}
	}

	tasks, err := r.taskProvider.List()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.EnvironmentID == environmentID {
			if err := r.deleteTask(jobID, task.TaskID); err != nil {
				return err
			}
		}
	}

	return catchAndRetry(15, func() (shouldRetry bool, err error) {
		if err := r.environmentProvider.Delete(environmentID); err != nil {
			switch err := err.(type) {
			case awserr.Error:
				return err.Code() == "DependencyViolation", err
			default:
				return false, err
			}
		}

		return false, nil
	})
}

func (r *JobRunner) deleteLoadBalancer(jobID, loadBalancerID string) error {
	return r.loadBalancerProvider.Delete(loadBalancerID)
}

func (r *JobRunner) deleteService(jobID, serviceID string) error {
	return r.serviceProvider.Delete(serviceID)
}

func (r *JobRunner) deleteTask(jobID, taskID string) error {
	return r.taskProvider.Delete(taskID)
}

func (r *JobRunner) updateEnvironment(jobID, request string) error {
	var req models.UpdateEnvironmentRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return errors.New(errors.InvalidRequest, err)
	}

	return r.environmentProvider.Update(req)
}

func (r *JobRunner) updateLoadBalancer(jobID, request string) error {
	var req models.UpdateLoadBalancerRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return errors.New(errors.InvalidRequest, err)
	}

	return r.loadBalancerProvider.Update(req)
}

func (r *JobRunner) updateService(jobID, request string) error {
	var req models.UpdateServiceRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return errors.New(errors.InvalidRequest, err)
	}

	return r.serviceProvider.Update(req)
}

func catchAndRetry(max int, fn func() (shouldRetry bool, err error)) error {
	var shouldRetry bool
	var err error

	for i := 0; i < max; i++ {
		shouldRetry, err = fn()
		if !shouldRetry {
			break
		}

		time.Sleep(time.Second * 5)
	}

	return err
}
