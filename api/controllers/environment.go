package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/entity"
	"github.com/quintilesims/layer0/api/scheduler"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/job"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type EnvironmentController struct {
	EnvironmentProvider     entity.EnvironmentProvider
	JobScheduler scheduler.JobScheduler
}

func NewEnvironmentController(p entity.Provider, j scheduler.JobScheduler) *EnvironmentController {
	return &EnvironmentController{
		Provider:     p,
		JobScheduler: j,
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

	environment := e.EnvironmentProvider.GetEnvironment("")
	if err := environment.Create(req); err != nil {
		return nil, err
	}

	environmentModel, err := environment.Model()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(202, environmentModel)
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
	environment := e.EnvironmentProvider.GetEnvironment(id)
	environmentModel, err := environment.Model()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, environmentModel)
}

func (e *EnvironmentController) ListEnvironments(c *fireball.Context) (fireball.Response, error) {
	environmentIDs, err := e.EnvironmentProvider.ListEnvironmentIDs()
	if err != nil {
		return nil, err
	}

	environmentSummaries := make([]*models.EnvironmentSummary, len(environmentIDs))
	for i, environmentID := range environmentIDs {
		environment := e.EnvironmentProvider.GetEnvironment(environmentID)
		environmentSummary, err := environment.Summary()
		if err != nil {
			return nil, err
		}

		environmentSummaries[i] = environmentSummary
	}

	return fireball.NewJSONResponse(200, environmentSummaries)

}
