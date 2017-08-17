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

type EnvironmentController struct {
	EnvironmentProvider provider.EnvironmentProvider
	JobScheduler        scheduler.JobScheduler
}

func NewEnvironmentController(e provider.EnvironmentProvider, j scheduler.JobScheduler) *EnvironmentController {
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

	model, err := e.EnvironmentProvider.Create(req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(202, model)
}

func (e *EnvironmentController) DeleteEnvironment(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	job := models.CreateJobRequest{
		JobType: job.DeleteEnvironmentJob,
		Request: id,
	}

	jobID, err := e.JobScheduler.ScheduleJob(job)
	if err != nil {
		return nil, err
	}

	return newJobResponse(jobID), nil
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
