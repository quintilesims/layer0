package models

import "fmt"

type CreateServiceRequest struct {
	DeployID       string `json:"deploy_id"`
	EnvironmentID  string `json:"environment_id"`
	LoadBalancerID string `json:"load_balancer_id"`
	ServiceName    string `json:"service_name"`
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

	return nil
}
