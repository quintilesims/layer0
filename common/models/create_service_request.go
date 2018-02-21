package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateServiceRequest struct {
	ServiceName    string `json:"service_name"`
	EnvironmentID  string `json:"environment_id"`
	DeployID       string `json:"deploy_id"`
	LoadBalancerID string `json:"load_balancer_id"`
	Scale          int    `json:"scale"`
}

func (c CreateServiceRequest) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("ServiceName is required")
	}

	if c.EnvironmentID == "" {
		return fmt.Errorf("EnvironmmentID is required")
	}

	if c.DeployID == "" {
		return fmt.Errorf("DeployID is required")
	}

	if c.Scale < 0 {
		return fmt.Errorf("Scale must be a positive integer")
	}

	return nil
}

func (s CreateServiceRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"deploy_id":        swagger.NewStringProperty(),
			"environment_id":   swagger.NewStringProperty(),
			"load_balancer_id": swagger.NewStringProperty(),
			"service_name":     swagger.NewStringProperty(),
			"scale":            swagger.NewIntProperty(),
		},
	}
}
