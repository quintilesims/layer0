package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type Service struct {
	Deployments      []Deployment `json:"deployments"`
	DesiredCount     int          `json:"desired_count"`
	EnvironmentID    string       `json:"environment_id"`
	EnvironmentName  string       `json:"environment_name"`
	LoadBalancerID   string       `json:"load_balancer_id"`
	LoadBalancerName string       `json:"load_balancer_name"`
	PendingCount     int          `json:"pending_count"`
	RunningCount     int          `json:"running_count"`
	ServiceID        string       `json:"service_id"`
	ServiceName      string       `json:"service_name"`
	Stateful         bool         `json:"stateful"`
}

func (s Service) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"deployments":        swagger.NewObjectSliceProperty("Deployment"),
			"desired_count":      swagger.NewIntProperty(),
			"environment_id":     swagger.NewStringProperty(),
			"environment_name":   swagger.NewStringProperty(),
			"load_balancer_id":   swagger.NewStringProperty(),
			"load_balancer_name": swagger.NewStringProperty(),
			"pending_count":      swagger.NewIntProperty(),
			"running_count":      swagger.NewIntProperty(),
			"service_id":         swagger.NewStringProperty(),
			"service_name":       swagger.NewStringProperty(),
			"stateful":           swagger.NewBoolProperty(),
		},
	}
}
