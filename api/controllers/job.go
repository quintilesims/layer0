package controllers

import (
	"github.com/quintilesims/layer0/api/entity"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type JobController struct {
	Provider entity.Provider
}

func NewJobController(p entity.Provider) *JobController {
	return &JobController{
		Provider: p,
	}
}

func (e *JobController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/job",
			Handlers: fireball.Handlers{
				"GET": e.ListJobs,
			},
		},
		{
			Path: "/job/:id",
			Handlers: fireball.Handlers{
				"GET":    e.GetJob,
				"DELETE": e.DeleteJob,
			},
		},
	}
}

func (e *JobController) DeleteJob(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	job := e.Provider.GetJob(id)
	if err := job.Delete(); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (e *JobController) GetJob(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	job := e.Provider.GetJob(id)
	jobModel, err := job.Model()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, jobModel)
}

func (e *JobController) ListJobs(c *fireball.Context) (fireball.Response, error) {
	jobIDs, err := e.Provider.ListJobIDs()
	if err != nil {
		return nil, err
	}

	jobModels := make([]*models.Job, len(jobIDs))
	for i, jobID := range jobIDs {
		job := e.Provider.GetJob(jobID)
		jobModel, err := job.Model()
		if err != nil {
			return nil, err
		}

		jobModels[i] = jobModel
	}

	return fireball.NewJSONResponse(200, jobModels)

}
