package models

import (
	"time"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type Deployment struct {
	Created       time.Time `json:"created"`
	DeployID      string    `json:"deploy_id"`
	DeployName    string    `json:"deploy_name"`
	DeployVersion string    `json:"deploy_version"`
	DesiredCount  int64     `json:"desired_count"`
	PendingCount  int64     `json:"pending_count"`
	RunningCount  int64     `json:"running_count"`
	Status        string    `json:"status"`
	Updated       time.Time `json:"updated"`
	DeploymentID  string    `json:"deployment_id"`
}

func (u Deployment) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"created":        swagger.NewStringProperty(),
			"deploy_id":      swagger.NewStringProperty(),
			"deploy_name":    swagger.NewStringProperty(),
			"deploy_version": swagger.NewIntProperty(),
			"desired_count":  swagger.NewIntProperty(),
			"pending_count":  swagger.NewIntProperty(),
			"running_count":  swagger.NewIntProperty(),
			"status":         swagger.NewStringProperty(),
			"updated":        swagger.NewStringProperty(),
			"deployment_id":  swagger.NewStringProperty(),
		},
	}
}
