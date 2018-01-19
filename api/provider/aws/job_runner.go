package aws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	scaler               scaler.Scaler
	dispatcher           *scaler.Dispatcher
}

func NewJobRunner(
	d provider.DeployProvider,
	e provider.EnvironmentProvider,
	l provider.LoadBalancerProvider,
	s provider.ServiceProvider,
	t provider.TaskProvider,
	ss scaler.Scaler,
	sd *scaler.Dispatcher,
) *JobRunner {
	return &JobRunner{
		deployProvider:       d,
		environmentProvider:  e,
		loadBalancerProvider: l,
		serviceProvider:      s,
		taskProvider:         t,
		scaler:               ss,
		dispatcher:           sd,
	}
}

func (r *JobRunner) Run(j models.Job) (string, error) {

	switch j.Type {
	case models.CreateDeployJob:
		return r.createDeploy(j.JobID, j.Request)
	case models.CreateEnvironmentJob:
		return r.createEnvironment(j.JobID, j.Request)
	case models.CreateLoadBalancerJob:
		return r.createLoadBalancer(j.JobID, j.Request)
	case models.CreateServiceJob:
		return r.createService(j.JobID, j.Request)
	case models.CreateTaskJob:
		return r.createTask(j.JobID, j.Request)
	case models.DeleteDeployJob:
		return r.deleteDeploy(j.JobID, j.Request)
	case models.DeleteEnvironmentJob:
		return r.deleteEnvironment(j.JobID, j.Request)
	case models.DeleteLoadBalancerJob:
		return r.deleteLoadBalancer(j.JobID, j.Request)
	case models.DeleteServiceJob:
		return r.deleteService(j.JobID, j.Request)
	case models.DeleteTaskJob:
		return r.deleteTask(j.JobID, j.Request)
	case models.ScaleEnvironmentJob:
		return r.scaleEnvironment(j.JobID, j.Request)
	case models.UpdateEnvironmentJob:
		return r.updateEnvironment(j.JobID, j.Request)
	case models.UpdateLoadBalancerJob:
		return r.updateLoadBalancer(j.JobID, j.Request)
	case models.UpdateServiceJob:
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

	serviceID, err := catchAndRetry(time.Minute*5, func() (result string, err error, shouldRetry bool) {
		log.Printf("[DEBUG] [JobRunner] Creating service %s", req.ServiceName)
		serviceID, err := r.serviceProvider.Create(req)
		if err != nil {
			log.Printf("[DEBUG] [JobRunner] Failed to create service %s: %v", req.ServiceName, err)
			return "", err, true
		}

		return serviceID, nil, false
	})
	if err != nil {
		return "", err
	}

	// scale up after the service has been added into the ecs service scheduler
	r.dispatcher.Dispatch(req.EnvironmentID)
	return serviceID, nil
}

func (r *JobRunner) createTask(jobID, request string) (string, error) {
	var req models.CreateTaskRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	// scale up prior to creating the task so we have room in the cluster
	r.dispatcher.Dispatch(req.EnvironmentID)

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
	return "", r.environmentProvider.Delete(environmentID)
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

func (r *JobRunner) scaleEnvironment(jobID, request string) (string, error) {
	return "", r.scaler.Scale(request)
}

func (r *JobRunner) updateEnvironment(jobID, request string) (string, error) {
	var req models.UpdateEnvironmentRequestJob
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	return req.EnvironmentID, r.environmentProvider.Update(req.EnvironmentID, req.UpdateEnvironmentRequest)
}

func (r *JobRunner) updateLoadBalancer(jobID, request string) (string, error) {
	var req models.UpdateLoadBalancerRequestJob
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	return req.LoadBalancerID, r.loadBalancerProvider.Update(req.LoadBalancerID, req.UpdateLoadBalancerRequest)
}

func (r *JobRunner) updateService(jobID, request string) (string, error) {
	var req models.UpdateServiceRequestJob
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		return "", errors.New(errors.InvalidRequest, err)
	}

	service, err := r.serviceProvider.Read(req.ServiceID)
	if err != nil {
		return "", err
	}

	r.dispatcher.Dispatch(service.EnvironmentID)
	return req.ServiceID, r.serviceProvider.Update(req.ServiceID, req.UpdateServiceRequest)
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
