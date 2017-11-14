package models

import (
	swagger "github.com/zpatrick/go-plugin-swagger"
)

type UpdateLoadBalancerRequestJob struct {
	LoadBalancerID string
	UpdateLoadBalancerRequest
}

type UpdateLoadBalancerRequest struct {
	Ports       *[]Port      `json:"ports"`
	HealthCheck *HealthCheck `json:"health_check"`
}

func (u UpdateLoadBalancerRequest) Validate() error {
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
			"ports":        swagger.NewObjectSliceProperty("Port"),
			"health_check": swagger.NewObjectProperty("HealthCheck"),
		},
	}
}
