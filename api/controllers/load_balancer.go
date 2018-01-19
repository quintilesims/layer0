package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type LoadBalancerController struct {
	LoadBalancerProvider provider.LoadBalancerProvider
}

func NewLoadBalancerController(l provider.LoadBalancerProvider) *LoadBalancerController {
	return &LoadBalancerController{
		LoadBalancerProvider: l,
	}
}

func (l *LoadBalancerController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/loadbalancer",
			Handlers: fireball.Handlers{
				"GET":  l.listLoadBalancers,
				"POST": l.createLoadBalancer,
			},
		},
		{
			Path: "/loadbalancer/:id",
			Handlers: fireball.Handlers{
				"GET":    l.readLoadBalancer,
				"DELETE": l.deleteLoadBalancer,
				"PATCH":  l.updateLoadBalancer,
			},
		},
	}
}

func (l *LoadBalancerController) createLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateLoadBalancerRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	loadBalancerID, err := l.LoadBalancerProvider.Create(req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, models.CreateEntityResponse{EntityID: loadBalancerID})
}

func (l *LoadBalancerController) deleteLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	loadBalancerID := c.PathVariables["id"]
	if err := l.LoadBalancerProvider.Delete(loadBalancerID); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (l *LoadBalancerController) listLoadBalancers(c *fireball.Context) (fireball.Response, error) {
	loadBalancers, err := l.LoadBalancerProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, loadBalancers)
}

func (l *LoadBalancerController) readLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	loadBalancerID := c.PathVariables["id"]
	loadBalancer, err := l.LoadBalancerProvider.Read(loadBalancerID)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, loadBalancer)
}

func (l *LoadBalancerController) updateLoadBalancer(c *fireball.Context) (fireball.Response, error) {
	loadBalancerID := c.PathVariables["id"]

	var req models.UpdateLoadBalancerRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := l.LoadBalancerProvider.Update(loadBalancerID, req); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}
