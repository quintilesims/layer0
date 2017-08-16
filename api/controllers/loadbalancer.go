package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/api/scheduler"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/job"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type LoadBalancerController struct {
	LoadBalancerProvider provider.LoadBalancerProvider
	JobScheduler         scheduler.JobScheduler
}

func NewLoadBalancerController(e provider.LoadBalancerProvider, j scheduler.JobScheduler) *LoadBalancerController {
	return &LoadBalancerController{
		LoadBalancerProvider: e,
		JobScheduler:         j,
	}
}

func (e *LoadBalancerController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/loadbalancer",
			Handlers: fireball.Handlers{
				"GET":  e.ListLoadBalancers,
				"POST": e.CreateLoadBalancer,
			},
		},
		{
			Path: "/loadbalancer/:id",
			Handlers: fireball.Handlers{
				"GET":    e.GetLoadBalancer,
				"DELETE": e.DeleteLoadBalancer,
			},
		},
	}
}

func (e *LoadBalancerController) CreateLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateLoadBalancerRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	model, err := e.LoadBalancerProvider.Create(req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(202, model)
}

func (e *LoadBalancerController) DeleteLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	job := models.CreateJobRequest{
		JobType: job.DeleteLoadBalancerJob,
		Request: id,
	}

	jobID, err := e.JobScheduler.ScheduleJob(job)
	if err != nil {
		return nil, err
	}

	return newJobResponse(jobID), nil
}

func (e *LoadBalancerController) GetLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	model, err := e.LoadBalancerProvider.Read(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, model)
}

func (e *LoadBalancerController) ListLoadBalancers(c *fireball.Context) (fireball.Response, error) {
	summaries, err := e.LoadBalancerProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, summaries)

}
