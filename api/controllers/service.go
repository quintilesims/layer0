package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/api/scheduler"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type ServiceController struct {
	ServiceProvider provider.ServiceProvider
	JobScheduler    scheduler.JobScheduler
}

func NewServiceController(s provider.ServiceProvider) *ServiceController {
	return &ServiceController{
		ServiceProvider: s,
	}
}

func (s *ServiceController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/service",
			Handlers: fireball.Handlers{
				"GET":  s.ListServices,
				"POST": s.CreateService,
			},
		},
		{
			Path: "/service/:id",
			Handlers: fireball.Handlers{
				"GET":    s.GetService,
				"DELETE": s.DeleteService,
			},
		},
	}
}

func (s *ServiceController) CreateService(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateServiceRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	model, err := s.ServiceProvider.Create(req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(202, model)
}

func (s *ServiceController) DeleteService(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	if err := s.ServiceProvider.Delete(id); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (s *ServiceController) GetService(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	model, err := s.ServiceProvider.Read(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, model)
}

func (s *ServiceController) ListServices(c *fireball.Context) (fireball.Response, error) {
	summaries, err := s.ServiceProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, summaries)

}
