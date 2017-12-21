package controllers

import (
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
	"github.com/zpatrick/fireball"
)

type AdminController struct {
	Context *cli.Context
	Version string
}

func NewAdminController(c *cli.Context, version string) *AdminController {
	return &AdminController{
		Context: c,
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
		{
			Path: "/admin/health",
			Handlers: fireball.Handlers{
				"GET": a.GetHealth,
			},
		},
	}
}

func (a *AdminController) GetConfig(c *fireball.Context) (fireball.Response, error) {
	model := models.APIConfig{
		Instance:       a.Context.String(config.FlagInstance.GetName()),
		VPCID:          a.Context.String(config.FlagAWSVPC.GetName()),
		Version:        a.Version,
		PublicSubnets:  a.Context.StringSlice(config.FlagAWSPublicSubnets.GetName()),
		PrivateSubnets: a.Context.StringSlice(config.FlagAWSPrivateSubnets.GetName()),
	}

	return fireball.NewJSONResponse(200, model)
}

func (a *AdminController) GetHealth(c *fireball.Context) (fireball.Response, error) {
	return fireball.NewJSONResponse(200, "")
}
