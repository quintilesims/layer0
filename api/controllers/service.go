package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type ServiceController struct {
	ServiceProvider provider.ServiceProvider
	JobScheduler    job.Store
}

func NewServiceController(s provider.ServiceProvider, j job.Store) *ServiceController {
	return &ServiceController{
		ServiceProvider: s,
		JobScheduler:    j,
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

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	return scheduleJob(s.JobScheduler, job.CreateServiceJob, req)
}

func (s *ServiceController) DeleteService(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	return scheduleJob(s.JobScheduler, job.DeleteServiceJob, id)
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
