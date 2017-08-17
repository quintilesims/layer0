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

func NewLoadBalancerController(l provider.LoadBalancerProvider, j scheduler.JobScheduler) *LoadBalancerController {
	return &LoadBalancerController{
		LoadBalancerProvider: l,
		JobScheduler:         j,
	}
}

func (l *LoadBalancerController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/loadbalancer",
			Handlers: fireball.Handlers{
				"GET":  l.ListLoadBalancers,
				"POST": l.CreateLoadBalancer,
			},
		},
		{
			Path: "/loadbalancer/:id",
			Handlers: fireball.Handlers{
				"GET":    l.GetLoadBalancer,
				"DELETE": l.DeleteLoadBalancer,
			},
		},
	}
}

func (l *LoadBalancerController) CreateLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateLoadBalancerRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
                return nil, errors.New(errors.InvalidRequest, err)
        }

	model, err := l.LoadBalancerProvider.Create(req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(202, model)
}

func (l *LoadBalancerController) DeleteLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	job := models.CreateJobRequest{
		JobType: job.DeleteLoadBalancerJob,
		Request: id,
	}

	jobID, err := l.JobScheduler.ScheduleJob(job)
	if err != nil {
		return nil, err
	}

	return newJobResponse(jobID), nil
}

func (l *LoadBalancerController) GetLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	model, err := l.LoadBalancerProvider.Read(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, model)
}

func (l *LoadBalancerController) ListLoadBalancers(c *fireball.Context) (fireball.Response, error) {
	summaries, err := l.LoadBalancerProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, summaries)

}
