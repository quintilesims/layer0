package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type Task struct {
	TaskID          string `json:"task_id"`
	TaskName        string `json:"task_name"`
	EnvironmentID   string `json:"environment_id"`
	EnvironmentName string `json:"environment_name"`
	DeployID        string `json:"deploy_id"`
	DeployName      string `json:"deploy_name"`
	DeployVersion   int    `json:"deploy_version"`

	// todo: fix copies field
	//Copies       []TaskCopy `json:"copies"`
	DesiredCount int `json:"desired_count"`
	PendingCount int `json:"pending_count"`
	RunningCount int `json:"running_count"`
}

func (t Task) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"task_id":          swagger.NewStringProperty(),
			"task_name":        swagger.NewStringProperty(),
			"environment_id":   swagger.NewStringProperty(),
			"environment_name": swagger.NewStringProperty(),
			"deploy_id":        swagger.NewStringProperty(),
			"deploy_name":      swagger.NewStringProperty(),
			"deploy_version":   swagger.NewIntProperty(),

			// todo: copies
			"desired_count": swagger.NewIntProperty(),
			"pending_count": swagger.NewIntProperty(),
			"running_count": swagger.NewIntProperty(),
		},
	}
}
