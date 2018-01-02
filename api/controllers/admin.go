package controllers

import (
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type AdminController struct {
	Config  config.APIConfig
	Version string
}

func NewAdminController(c config.APIConfig, version string) *AdminController {
	return &AdminController{
		Config:  c,
		Version: version,
	}
}

func (a *AdminController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/admin/config",
			Handlers: fireball.Handlers{
				"GET": a.GetConfig,
			},
		},
	}
}

func (a *AdminController) GetConfig(c *fireball.Context) (fireball.Response, error) {
	model := models.APIConfig{
		Instance:       a.Config.Instance(),
		VPCID:          a.Config.VPC(),
		Version:        a.Version,
		PublicSubnets:  a.Config.PublicSubnets(),
		PrivateSubnets: a.Config.PrivateSubnets(),
	}

	return fireball.NewJSONResponse(200, model)
}
