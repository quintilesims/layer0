package models

import "fmt"

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
