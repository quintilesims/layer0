package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type DeployController struct {
	DeployProvider provider.DeployProvider
}

func NewDeployController(d provider.DeployProvider) *DeployController {
	return &DeployController{
		DeployProvider: d,
	}
}

func (d *DeployController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/deploy",
			Handlers: fireball.Handlers{
				"GET":  d.listDeploys,
				"POST": d.createDeploy,
			},
		},
		{
			Path: "/deploy/:id",
			Handlers: fireball.Handlers{
				"GET":    d.readDeploy,
				"DELETE": d.deleteDeploy,
			},
		},
	}
}

func (d *DeployController) createDeploy(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateDeployRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	deployID, err := d.DeployProvider.Create(req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, models.CreateEntityResponse{EntityID: deployID})
}

func (d *DeployController) deleteDeploy(c *fireball.Context) (fireball.Response, error) {
	deployID := c.PathVariables["id"]
	if err := d.DeployProvider.Delete(deployID); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (d *DeployController) listDeploys(c *fireball.Context) (fireball.Response, error) {
	deploys, err := d.DeployProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, deploys)
}

func (d *DeployController) readDeploy(c *fireball.Context) (fireball.Response, error) {
	deployID := c.PathVariables["id"]
	deploy, err := d.DeployProvider.Read(deployID)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, deploy)
}
