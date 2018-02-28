package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type LoadBalancer struct {
	EnvironmentID    string      `json:"environment_id"`
	EnvironmentName  string      `json:"environment_name"`
	HealthCheck      HealthCheck `json:"health_check"`
	IdleTimeout      int         `json:"idle_timeout"`
	IsPublic         bool        `json:"is_public"`
	LoadBalancerID   string      `json:"load_balancer_id"`
	LoadBalancerName string      `json:"load_balancer_name"`
	Ports            []Port      `json:"ports"`
	ServiceID        string      `json:"service_id"`
	ServiceName      string      `json:"service_name"`
	URL              string      `json:"url"`
}

func (l LoadBalancer) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"environment_id":     swagger.NewStringProperty(),
			"environment_name":   swagger.NewStringProperty(),
			"health_check":       swagger.NewObjectProperty("HealthCheck"),
			"idle_timeout":       swagger.NewIntProperty(),
			"is_public":          swagger.NewBoolProperty(),
			"load_balancer_id":   swagger.NewStringProperty(),
			"load_balancer_name": swagger.NewStringProperty(),
			"ports":              swagger.NewObjectSliceProperty("Port"),
			"service_id":         swagger.NewStringProperty(),
			"service_name":       swagger.NewStringProperty(),
			"url":                swagger.NewStringProperty(),
		},
	}
}
