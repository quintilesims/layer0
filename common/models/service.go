package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type Service struct {
	Deployments      []Deployment `json:"deployments"`
	DesiredCount     int64        `json:"desired_count"`
	EnvironmentID    string       `json:"environment_id"`
	EnvironmentName  string       `json:"environment_name"`
	LoadBalancerID   string       `json:"load_balancer_id"`
	LoadBalancerName string       `json:"load_balancer_name"`
	PendingCount     int64        `json:"pending_count"`
	RunningCount     int64        `json:"running_count"`
	ServiceID        string       `json:"service_id"`
	ServiceName      string       `json:"service_name"`
}

func (s Service) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"deployments":        swagger.NewObjectSliceProperty("deployments"),
			"desired_count":      swagger.NewIntProperty(),
			"environment_id":     swagger.NewStringProperty(),
			"environment_name":   swagger.NewStringProperty(),
			"load_balancer_id":   swagger.NewStringProperty(),
			"load_balancer_name": swagger.NewStringProperty(),
			"pending_count":      swagger.NewIntProperty(),
			"running_count":      swagger.NewIntProperty(),
			"service_id":         swagger.NewStringProperty(),
			"service_name":       swagger.NewStringProperty(),
		},
	}
}
