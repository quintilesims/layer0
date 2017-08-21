package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type HealthCheck struct {
	Target             string `json:"target"`
	Interval           int    `json:"interval"`
	Timeout            int    `json:"timeout"`
	HealthyThreshold   int    `json:"healthy_threshold"`
	UnhealthyThreshold int    `json:"unhealthy_threshold"`
}

func (h HealthCheck) Validate() error {
	if h.Target == "" {
		return fmt.Errorf("Target is required")
	}

	if h.Interval == 0 {
		return fmt.Errorf("Interval is required")
	}

	if h.Timeout == 0 {
		return fmt.Errorf("Timeout is required")
	}

	if h.HealthyThreshold == 0 {
		return fmt.Errorf("HealthyThreshold is required")
	}

	if h.UnhealthyThreshold == 0 {
		return fmt.Errorf("UnhealthyThreshold is required")
	}

	return nil
}

func (h HealthCheck) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"target":              swagger.NewStringProperty(),
			"interval":            swagger.NewIntProperty(),
			"timeout":             swagger.NewIntProperty(),
			"healthy_threshold":   swagger.NewIntProperty(),
			"unhealthy_threshold": swagger.NewIntProperty(),
		},
	}
}
