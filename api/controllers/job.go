package controllers

import (
	"github.com/quintilesims/layer0/api/job"
	"github.com/zpatrick/fireball"
)

type JobController struct {
	JobStore job.Store
}

func NewJobController(j job.Store) *JobController {
	return &JobController{
		JobStore: j,
	}
}

func (j *JobController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/job",
			Handlers: fireball.Handlers{
				"GET": j.ListJobs,
			},
		},
		{
			Path: "/job/:id",
			Handlers: fireball.Handlers{
				"GET":    j.GetJob,
				"DELETE": j.DeleteJob,
			},
		},
	}
}

func (j *JobController) DeleteJob(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	if err := j.JobStore.Delete(id); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (j *JobController) GetJob(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	job, err := j.JobStore.SelectByID(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, job)
}

func (j *JobController) ListJobs(c *fireball.Context) (fireball.Response, error) {
	jobs, err := j.JobStore.SelectAll()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, jobs)

}
