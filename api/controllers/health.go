package controllers

import (
	"github.com/quintilesims/layer0/common/config"
	"github.com/zpatrick/fireball"
)

type HealthController struct {
	Config  config.APIConfig
	Version string
}

func NewHealthController(c config.APIConfig, version string) *HealthController {
	return &HealthController{
		Config:  c,
		Version: version,
	}
}

func (a *HealthController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/health",
			Handlers: fireball.Handlers{
				"GET": a.GetHealth,
			},
		},
	}
}

func (a *HealthController) GetHealth(c *fireball.Context) (fireball.Response, error) {
	return fireball.NewJSONResponse(200, "")
}
