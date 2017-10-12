package aws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
}

func NewJobRunner(
	d provider.DeployProvider,
	e provider.EnvironmentProvider,
	l provider.LoadBalancerProvider,
	s provider.ServiceProvider,
	t provider.TaskProvider,
	scaler *scaler.Dispatcher,
) *JobRunner {
	return &JobRunner{
		deployProvider:       d,
		environmentProvider:  e,
		loadBalancerProvider: l,
		serviceProvider:      s,
		taskProvider:         t,
		scaler:               scaler,
	}
}

func (r *JobRunner) Run(j models.Job) (string, error) {
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
		return "", fmt.Errorf("Unrecognized JobType '%s'", j.Type)
	}
}

func (r *JobRunner) createDeploy(jobID, request string) (string, error) {
	var req models.CreateDeployRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	deployID, err := r.deployProvider.Create(req)
	if err != nil {
		return "", err
	}

	return deployID, nil
}

func (r *JobRunner) createEnvironment(jobID, request string) (string, error) {
	var req models.CreateEnvironmentRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	environmentID, err := r.environmentProvider.Create(req)
	if err != nil {
		return "", err
	}

	return environmentID, nil
}

func (r *JobRunner) createLoadBalancer(jobID, request string) (string, error) {
	var req models.CreateLoadBalancerRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	loadBalancerID, err := r.loadBalancerProvider.Create(req)
	if err != nil {
		return "", err
	}

	return loadBalancerID, nil
}

func (r *JobRunner) createService(jobID, request string) (string, error) {
	var req models.CreateServiceRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	serviceID, err := r.serviceProvider.Create(req)
	if err != nil {
		return "", err
	}

	// scale up after the service has been added into the ecs service scheduler
	r.scaler.ScheduleRun(req.EnvironmentID)
	return serviceID, nil
}

func (r *JobRunner) createTask(jobID, request string) (string, error) {
	var req models.CreateTaskRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	// scale up prior to creating the task so we have room in the cluster
	r.scaler.ScheduleRun(req.EnvironmentID)

	return catchAndRetry(time.Hour*24, func() (result string, err error, shouldRetry bool) {
		log.Printf("[DEBUG] [JobRunner] Creating task %s", req.TaskName)
		taskID, err := r.taskProvider.Create(req)
		if err != nil {
			log.Printf("[DEBUG] [JobRunner] Failed to create task %s: %v", req.TaskName, err)
			return "", err, true
		}

		return taskID, nil, false
	})
}

func (r *JobRunner) deleteDeploy(jobID, deployID string) (string, error) {
	return "", r.deployProvider.Delete(deployID)
}

func (r *JobRunner) deleteEnvironment(jobID, environmentID string) (string, error) {
	loadBalancers, err := r.loadBalancerProvider.List()
	if err != nil {
		return "", err
	}

	for _, loadBalancer := range loadBalancers {
		if loadBalancer.EnvironmentID == environmentID {
			if _, err := r.deleteLoadBalancer(jobID, loadBalancer.LoadBalancerID); err != nil {
				return "", err
			}
		}
	}

	services, err := r.serviceProvider.List()
	if err != nil {
		return "", err
	}

	for _, service := range services {
		if service.EnvironmentID == environmentID {
			if _, err := r.deleteService(jobID, service.ServiceID); err != nil {
				return "", err
			}
		}
	}

	tasks, err := r.taskProvider.List()
	if err != nil {
		return "", err
	}

	for _, task := range tasks {
		if task.EnvironmentID == environmentID {
			if _, err := r.deleteTask(jobID, task.TaskID); err != nil {
				return "", err
			}
		}
	}

	return catchAndRetry(time.Minute*15, func() (result string, err error, shouldRetry bool) {
		log.Printf("[DEBUG] [JobRunner] Deleting environment %s", environmentID)
		if err := r.environmentProvider.Delete(environmentID); err != nil {
			log.Printf("[DEBUG] [JobRunner] Failed to delete environment %s: %v", environmentID, err)
			return "", err, true
		}

		return "", nil, false
	})
}

func (r *JobRunner) deleteLoadBalancer(jobID, loadBalancerID string) (string, error) {
	return catchAndRetry(time.Minute*15, func() (result string, err error, shouldRetry bool) {
		log.Printf("[DEBUG] [JobRunner] Deleting load balancer %s", loadBalancerID)
		if err := r.loadBalancerProvider.Delete(loadBalancerID); err != nil {
			log.Printf("[DEBUG] [JobRunner] Failed to delete load balancer %s: %v", loadBalancerID, err)
			return "", err, true
		}

		return "", nil, false
	})
}

func (r *JobRunner) deleteService(jobID, serviceID string) (string, error) {
	return "", r.serviceProvider.Delete(serviceID)
}

func (r *JobRunner) deleteTask(jobID, taskID string) (string, error) {
	return "", r.taskProvider.Delete(taskID)
}

func (r *JobRunner) updateEnvironment(jobID, request string) (string, error) {
	var req models.UpdateEnvironmentRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	return req.EnvironmentID, r.environmentProvider.Update(req)
}

func (r *JobRunner) updateLoadBalancer(jobID, request string) (string, error) {
	var req models.UpdateLoadBalancerRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	return req.LoadBalancerID, r.loadBalancerProvider.Update(req)
}

func (r *JobRunner) updateService(jobID, request string) (string, error) {
	var req models.UpdateServiceRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	return req.ServiceID, r.serviceProvider.Update(req)
}

func catchAndRetry(timeout time.Duration, fn func() (result string, err error, shouldRetry bool)) (string, error) {
	var result string
	var err error
	var shouldRetry bool

	for start := time.Now(); time.Since(start) < timeout; {
		result, err, shouldRetry = fn()
		if !shouldRetry {
			break
		}

		time.Sleep(time.Second * 5)
	}

	return result, err
}
