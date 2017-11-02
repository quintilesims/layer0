package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateServiceRequest struct {
	DeployID       string `json:"deploy_id"`
	EnvironmentID  string `json:"environment_id"`
	LoadBalancerID string `json:"load_balancer_id"`
	ServiceName    string `json:"service_name"`
	Scale          int    `json:"scale"`
}

func (c CreateServiceRequest) Validate() error {
	if c.DeployID == "" {
		return fmt.Errorf("DeployID is required")
	}

	if c.EnvironmentID == "" {
		return fmt.Errorf("EnvironmmentID is required")
	}

	if c.ServiceName == "" {
		return fmt.Errorf("ServiceName is required")
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
