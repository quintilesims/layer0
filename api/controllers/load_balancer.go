package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type LoadBalancerController struct {
	LoadBalancerProvider provider.LoadBalancerProvider
	JobStore             job.Store
	TagStore             tag.Store
}

func NewLoadBalancerController(l provider.LoadBalancerProvider, j job.Store, t tag.Store) *LoadBalancerController {
	return &LoadBalancerController{
		LoadBalancerProvider: l,
		JobStore:             j,
		TagStore:             t,
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
				"GET":    l.ReadLoadBalancer,
				"PATCH":  l.UpdateLoadBalancer,
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

	return createJob(l.TagStore, l.JobStore, job.CreateLoadBalancerJob, req)
}

func (l *LoadBalancerController) DeleteLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	return createJob(l.TagStore, l.JobStore, job.DeleteLoadBalancerJob, id)
}

func (l *LoadBalancerController) ReadLoadBalancer(c *fireball.Context) (fireball.Response, error) {
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

func (l *LoadBalancerController) UpdateLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	var req models.UpdateLoadBalancerRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	jobRequest := models.UpdateLoadBalancerRequestJob{id, req}
	return createJob(l.TagStore, l.JobStore, job.UpdateLoadBalancerJob, jobRequest)
}
