package controllers

import (
	"os"
	"strings"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type AdminController struct {
	JobStore job.Store
}

func NewAdminController(j job.Store) *AdminController {
	return &AdminController{
		JobStore: j,
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
	publicSubnets := []string{}
	privateSubnets := []string{}

	VPCID := os.Getenv(config.ENVVAR_AWS_VPC)
	prefix := os.Getenv("LAYER0_PREFIX")

	subnet := os.Getenv(config.ENVVAR_AWS_PRIVATE_SUBNETS)
	for _, subnet := range strings.Split(subnet, ",") {
		publicSubnets = append(privateSubnets, subnet)
	}

	subnet = os.Getenv(config.ENVVAR_AWS_PUBLIC_SUBNETS)
	for _, subnet := range strings.Split(subnet, ",") {
		publicSubnets = append(publicSubnets, subnet)
	}

	model := models.APIConfig{
		Prefix:         prefix,
		VPCID:          VPCID,
		PublicSubnets:  publicSubnets,
		PrivateSubnets: publicSubnets,
	}

	return fireball.NewJSONResponse(200, model)
}

func (a *AdminController) GetHealth(c *fireball.Context) (fireball.Response, error) {
	return fireball.NewJSONResponse(200, "")
}
