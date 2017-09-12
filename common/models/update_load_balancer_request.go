package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type UpdateLoadBalancerRequest struct {
	LoadBalancerID string
	Ports          *[]Port     `json:"ports"`
	HealthCheck    HealthCheck `json:"health_check"`
}

func (u UpdateLoadBalancerRequest) Validate() error {
	if u.LoadBalancerID == "" {
		return fmt.Errorf("LoadBalancerID is required")
	}

	return nil
}

func (u UpdateLoadBalancerRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"target":              swagger.NewStringProperty(),
			"interval":            swagger.NewIntProperty(),
			"timeout":             swagger.NewIntProperty(),
			"healthy_threshold":   swagger.NewIntProperty(),
			"unhealthy_threshold": swagger.NewIntProperty(),
			"certificate_name":    swagger.NewStringProperty(),
			"container_port":      swagger.NewIntProperty(),
			"host_port":           swagger.NewIntProperty(),
			"protocol":            swagger.NewStringProperty(),
		},
	}
}
