package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type Task struct {
	Containers      []Container `json:"containers"`
	DeployID        string      `json:"deploy_id"`
	DeployName      string      `json:"deploy_name"`
	DeployVersion   string      `json:"deploy_version"`
	EnvironmentID   string      `json:"environment_id"`
	EnvironmentName string      `json:"environment_name"`
	Status          string      `json:"status"`
	TaskID          string      `json:"task_id"`
	TaskName        string      `json:"task_name"`
	Stateful        bool        `json:"stateful"`
}

func (t Task) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"containers":       swagger.NewObjectSliceProperty("Container"),
			"deploy_id":        swagger.NewStringProperty(),
			"deploy_name":      swagger.NewStringProperty(),
			"deploy_version":   swagger.NewStringProperty(),
			"environment_id":   swagger.NewStringProperty(),
			"environment_name": swagger.NewStringProperty(),
			"status":           swagger.NewStringProperty(),
			"task_id":          swagger.NewStringProperty(),
			"task_name":        swagger.NewStringProperty(),
			"stateful":         swagger.NewBoolProperty(),
		},
	}
}
