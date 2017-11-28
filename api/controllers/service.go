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

type ServiceController struct {
	ServiceProvider provider.ServiceProvider
	JobStore        job.Store
	TagStore        tag.Store
}

func NewServiceController(s provider.ServiceProvider, j job.Store, t tag.Store) *ServiceController {
	return &ServiceController{
		ServiceProvider: s,
		JobStore:        j,
		TagStore:        t,
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
				"PATCH":  s.UpdateService,
				"DELETE": s.DeleteService,
			},
		},
		{
			Path: "/service/:id/logs",
			Handlers: fireball.Handlers{
				"GET": s.GetServiceLogs,
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

	return createJob(s.JobStore, models.CreateServiceJob, req)
}

func (s *ServiceController) DeleteService(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	return createJob(s.JobStore, models.DeleteServiceJob, id)
}

func (s *ServiceController) GetService(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	model, err := s.ServiceProvider.Read(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, model)
}

func (s *ServiceController) GetServiceLogs(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	tail, start, end, err := parseLoggingQuery(c.Request.URL.Query())
	if err != nil {
		return nil, err
	}

	logs, err := s.ServiceProvider.Logs(id, tail, start, end)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, logs)
}

func (s *ServiceController) ListServices(c *fireball.Context) (fireball.Response, error) {
	summaries, err := s.ServiceProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, summaries)

}

func (s *ServiceController) UpdateService(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	var req models.UpdateServiceRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	jobRequest := models.UpdateServiceRequestJob{id, req}
	return createJob(s.JobStore, models.UpdateServiceJob, jobRequest)
}
