package controllers

import (
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type AdminController struct {
	AdminProvider provider.AdminProvider
	Config        config.APIConfig
	Version       string
}

func NewAdminController(a provider.AdminProvider, c config.APIConfig, version string) *AdminController {
	return &AdminController{
		AdminProvider: a,
		Config:        c,
		Version:       version,
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
			Path: "/admin/logs",
			Handlers: fireball.Handlers{
				"GET": a.readLogs,
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

func (a *AdminController) readLogs(c *fireball.Context) (fireball.Response, error) {
	tail, start, end, err := ParseLoggingQuery(c.Request.URL.Query())
	if err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	logs, err := a.AdminProvider.Logs(tail, start, end)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, logs)
}
