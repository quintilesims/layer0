package models

import (
	"fmt"
	"strings"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

const (
	ClassicLoadBalancerType     = "clb"
	ApplicationLoadBalancerType = "alb"
)

type CreateLoadBalancerRequest struct {
	LoadBalancerName string      `json:"load_balancer_name"`
	LoadBalancerType string      `json:"load_balancertype"`
	EnvironmentID    string      `json:"environment_id"`
	IsPublic         bool        `json:"is_public"`
	Ports            []Port      `json:"ports"`
	HealthCheck      HealthCheck `json:"health_check"`
}

func (c CreateLoadBalancerRequest) Validate() error {
	if c.LoadBalancerName == "" {
		return fmt.Errorf("LoadBalancerName is required")
	}

	if c.EnvironmentID == "" {
		return fmt.Errorf("Environment ID is required")
	}

	if !strings.EqualFold(c.LoadBalancerType, ClassicLoadBalancerType) &&
		!strings.EqualFold(c.LoadBalancerType, ApplicationLoadBalancerType) {
		return fmt.Errorf("%s is not a support load balancer type", c.LoadBalancerType)
	}

	for _, port := range c.Ports {
		if err := port.Validate(); err != nil {
			return err
		}
	}

	return c.HealthCheck.Validate()
}

func (l CreateLoadBalancerRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"load_balancer_name": swagger.NewStringProperty(),
			"load_balancer_type": swagger.NewStringProperty(),
			"environment_id":     swagger.NewStringProperty(),
			"is_public":          swagger.NewBoolProperty(),
			"ports":              swagger.NewObjectSliceProperty("Port"),
			"health_check":       swagger.NewObjectProperty("HealthCheck"),
		},
	}
}
