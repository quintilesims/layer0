package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type EnvironmentController struct {
	EnvironmentProvider provider.EnvironmentProvider
}

func NewEnvironmentController(e provider.EnvironmentProvider) *EnvironmentController {
	return &EnvironmentController{
		EnvironmentProvider: e,
	}
}

func (e *EnvironmentController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/environment",
			Handlers: fireball.Handlers{
				"GET":  e.listEnvironments,
				"POST": e.createEnvironment,
			},
		},
		{
			Path: "/environment/:id",
			Handlers: fireball.Handlers{
				"GET":    e.readEnvironment,
				"DELETE": e.deleteEnvironment,
			},
		},
		{
			Path: "/environment/:id",
			Handlers: fireball.Handlers{
				"PATCH": e.updateEnvironment,
			},
		},
		{
			Path: "/environment/:id/logs",
			Handlers: fireball.Handlers{
				"GET": e.readEnvironmentLogs,
			},
		},
	}
}

func (e *EnvironmentController) createEnvironment(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateEnvironmentRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	environmentID, err := e.EnvironmentProvider.Create(req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, models.CreateEntityResponse{EntityID: environmentID})
}

func (e *EnvironmentController) deleteEnvironment(c *fireball.Context) (fireball.Response, error) {
	environmentID := c.PathVariables["id"]
	if err := e.EnvironmentProvider.Delete(environmentID); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (e *EnvironmentController) listEnvironments(c *fireball.Context) (fireball.Response, error) {
	environments, err := e.EnvironmentProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, environments)
}

func (e *EnvironmentController) readEnvironment(c *fireball.Context) (fireball.Response, error) {
	environmentID := c.PathVariables["id"]
	environment, err := e.EnvironmentProvider.Read(environmentID)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, environment)
}

func (e *EnvironmentController) readEnvironmentLogs(c *fireball.Context) (fireball.Response, error) {
	environmentID := c.PathVariables["id"]
	tail, start, end, err := client.ParseLoggingQuery(c.Request.URL.Query())
	if err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	logs, err := e.EnvironmentProvider.Logs(environmentID, tail, start, end)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, logs)
}

func (e *EnvironmentController) updateEnvironment(c *fireball.Context) (fireball.Response, error) {
	environmentID := c.PathVariables["id"]

	var req models.UpdateEnvironmentRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := e.EnvironmentProvider.Update(environmentID, req); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}
