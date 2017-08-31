package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type DeployController struct {
	DeployProvider provider.DeployProvider
	JobScheduler   job.Scheduler
}

func NewDeployController(d provider.DeployProvider, j job.Scheduler) *DeployController {
	return &DeployController{
		DeployProvider: d,
		JobScheduler:   j,
	}
}

func (d *DeployController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/deploy",
			Handlers: fireball.Handlers{
				"GET":  d.ListDeploys,
				"POST": d.CreateDeploy,
			},
		},
		{
			Path: "/deploy/:id",
			Handlers: fireball.Handlers{
				"GET":    d.GetDeploy,
				"DELETE": d.DeleteDeploy,
			},
		},
	}
}

func (d *DeployController) CreateDeploy(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateDeployRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	return scheduleJob(d.JobScheduler, job.CreateDeployJob, req)
}

func (d *DeployController) DeleteDeploy(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	return scheduleJob(d.JobScheduler, job.DeleteDeployJob, id)
}

func (d *DeployController) GetDeploy(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	model, err := d.DeployProvider.Read(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, model)
}

func (d *DeployController) ListDeploys(c *fireball.Context) (fireball.Response, error) {
	summaries, err := d.DeployProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, summaries)

}
