package models

import (
	"fmt"
	"strings"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type HealthCheck struct {
	Target             string `json:"target"`
	Path               string `json:"path"`
	Interval           int    `json:"interval"`
	Timeout            int    `json:"timeout"`
	HealthyThreshold   int    `json:"healthy_threshold"`
	UnhealthyThreshold int    `json:"unhealthy_threshold"`
}

func (h HealthCheck) Validate() error {
	if h.Path != "" && !strings.HasPrefix(h.Path, "/") {
		return fmt.Errorf("expected healthcheck path '%s' to start with '/'", h.Path)
	}

	if h.Target != "" {
		split := strings.FieldsFunc(h.Target, func(r rune) bool {
			return r == ':' || r == '/'
		})

		protocol := strings.ToLower(split[0])
		if len(split) < 3 && (protocol == "https" || protocol == "http") {
			text := "HTTP & HTTPS targets must specify a port followed by a path.\n"
			text += "For example, HTTPS:443/health"
			return fmt.Errorf(text)
		}
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
			"path":                swagger.NewStringProperty(),
			"interval":            swagger.NewIntProperty(),
			"timeout":             swagger.NewIntProperty(),
			"healthy_threshold":   swagger.NewIntProperty(),
			"unhealthy_threshold": swagger.NewIntProperty(),
		},
	}
}
