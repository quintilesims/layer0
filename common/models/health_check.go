package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type HealthCheck struct {
	Target             string `json:"target"`
	Interval           int    `json:"interval"`
	Timeout            int    `json:"timeout"`
	HealthyThreshold   int    `json:"healthy_threshold"`
	UnhealthyThreshold int    `json:"unhealthy_threshold"`
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
