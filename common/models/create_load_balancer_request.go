package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateLoadBalancerRequest struct {
	LoadBalancerName string      `json:"load_balancer_name"`
	EnvironmentID    string      `json:"environment_id"`
	IsPublic         bool        `json:"is_public"`
	Ports            []Port      `json:"ports"`
	HealthCheck      HealthCheck `json:"health_check"`
}

func (c CreateLoadBalancerRequest) Validate() error {
	if c.LoadBalancerName == "" {
		return fmt.Errorf("LoadBalancer Name is required")
	}

	if c.EnvironmentID == "" {
		return fmt.Errorf("Environment ID is required")
	}

	return nil
}

func (l CreateLoadBalancerRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"load_balancer_name": swagger.NewStringProperty(),
			"environment_id":     swagger.NewStringProperty(),
			"is_public":          swagger.NewBoolProperty(),
			"ports":              swagger.NewObjectSliceProperty("Port"),
			"health_check":       swagger.NewObjectProperty("HealthCheck"),
		},
	}
}
