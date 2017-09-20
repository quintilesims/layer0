package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type UpdateLoadBalancerRequest struct {
	LoadBalancerID string       `json:"load_balancer_id"`
	Ports          *[]Port      `json:"ports"`
	HealthCheck    *HealthCheck `json:"health_check"`
}

func (u UpdateLoadBalancerRequest) Validate() error {
	if u.LoadBalancerID == "" {
		return fmt.Errorf("LoadBalancerID is required")
	}

	if u.Ports != nil {
		for _, port := range *u.Ports {
			if err := port.Validate(); err != nil {
				return err
			}
		}
	}

	if u.HealthCheck != nil {
		if err := u.HealthCheck.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (u UpdateLoadBalancerRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"load_balancer_id": swagger.NewStringProperty(),
			"ports":            swagger.NewObjectSliceProperty("Port"),
			"health_check":     swagger.NewObjectProperty("HealthCheck"),
		},
	}
}
