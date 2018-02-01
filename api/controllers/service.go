package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type ServiceController struct {
	ServiceProvider provider.ServiceProvider
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
				"GET":  s.listServices,
				"POST": s.createService,
			},
		},
		{
			Path: "/service/:id",
			Handlers: fireball.Handlers{
				"GET":    s.readService,
				"DELETE": s.deleteService,
				"PATCH":  s.updateService,
			},
		},
		{
			Path: "/service/:id/logs",
			Handlers: fireball.Handlers{
				"GET": s.readServiceLogs,
			},
		},
	}
}

func (s *ServiceController) createService(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateServiceRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	serviceID, err := s.ServiceProvider.Create(req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, models.CreateEntityResponse{EntityID: serviceID})
}

func (s *ServiceController) deleteService(c *fireball.Context) (fireball.Response, error) {
	serviceID := c.PathVariables["id"]
	if err := s.ServiceProvider.Delete(serviceID); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (s *ServiceController) listServices(c *fireball.Context) (fireball.Response, error) {
	services, err := s.ServiceProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, services)
}

func (s *ServiceController) readService(c *fireball.Context) (fireball.Response, error) {
	serviceID := c.PathVariables["id"]
	service, err := s.ServiceProvider.Read(serviceID)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, service)
}

func (s *ServiceController) readServiceLogs(c *fireball.Context) (fireball.Response, error) {
	serviceID := c.PathVariables["id"]
	tail, start, end, err := parseLoggingQuery(c.Request.URL.Query())
	if err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	logs, err := s.ServiceProvider.Logs(serviceID, tail, start, end)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, logs)
}

func (s *ServiceController) updateService(c *fireball.Context) (fireball.Response, error) {
	serviceID := c.PathVariables["id"]

	var req models.UpdateServiceRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := s.ServiceProvider.Update(serviceID, req); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}
