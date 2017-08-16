package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/api/scheduler"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type DeployController struct {
	DeployProvider provider.DeployProvider
	JobScheduler   scheduler.JobScheduler
}

func NewDeployController(p provider.DeployProvider) *DeployController {
	return &DeployController{
		DeployProvider: p,
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
				"DELETE": d.GetDeploy,
			},
		},
	}
}

func (e *DeployController) CreateDeploy(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateDeployRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	model, err := e.DeployProvider.Create(req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(202, model)
}

func (e *DeployController) DeleteDeploy(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	if err := e.DeployProvider.Delete(id); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (e *DeployController) GetDeploy(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	model, err := e.DeployProvider.Read(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, model)
}

func (e *DeployController) ListDeploys(c *fireball.Context) (fireball.Response, error) {
	summaries, err := e.DeployProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, summaries)

}
