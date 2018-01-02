package controllers

import (
	"github.com/zpatrick/fireball"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/health",
			Handlers: fireball.Handlers{
				"GET": h.GetHealth,
			},
		},
	}
}

func (h *HealthController) GetHealth(c *fireball.Context) (fireball.Response, error) {
	return fireball.NewJSONResponse(200, "")
}
