package controllers

import (
	"github.com/quintilesims/layer0/api/job"
	"github.com/zpatrick/fireball"
)

// todo: job test
type JobController struct {
	JobStore job.Store
}

func NewJobController(j job.Store) *JobController {
	return &JobController{
		JobStore: j,
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
	if err := e.JobStore.Delete(id); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (e *JobController) GetJob(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	job, err := e.JobStore.SelectByID(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, job)
}

func (e *JobController) ListJobs(c *fireball.Context) (fireball.Response, error) {
	jobs, err := e.JobStore.SelectAll()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, jobs)

}
