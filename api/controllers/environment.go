package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type EnvironmentController struct {
	EnvironmentProvider provider.EnvironmentProvider
	JobScheduler        job.Store
}

func NewEnvironmentController(e provider.EnvironmentProvider, j job.Store) *EnvironmentController {
	return &EnvironmentController{
		EnvironmentProvider: e,
		JobScheduler:        j,
	}
}

func (e *EnvironmentController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/environment",
			Handlers: fireball.Handlers{
				"GET":  e.ListEnvironments,
				"POST": e.CreateEnvironment,
			},
		},
		{
			Path: "/environment/:id",
			Handlers: fireball.Handlers{
				"GET":    e.GetEnvironment,
				"DELETE": e.DeleteEnvironment,
			},
		},
	}
}

func (e *EnvironmentController) CreateEnvironment(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateEnvironmentRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	return createJob(e.JobScheduler, job.CreateEnvironmentJob, req)
}

func (e *EnvironmentController) DeleteEnvironment(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	return createJob(e.JobScheduler, job.DeleteEnvironmentJob, id)
}

func (e *EnvironmentController) GetEnvironment(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	model, err := e.EnvironmentProvider.Read(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, model)
}

func (e *EnvironmentController) ListEnvironments(c *fireball.Context) (fireball.Response, error) {
	summaries, err := e.EnvironmentProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, summaries)

}
