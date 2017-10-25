package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateTaskRequest struct {
	ContainerOverrides []ContainerOverride `json:"container_overrides"`
	TaskName           string              `json:"task_name"`
	EnvironmentID      string              `json:"environment_id"`
	DeployID           string              `json:"deploy_id"`
}

func (c CreateTaskRequest) Validate() error {
	if c.TaskName == "" {
		return fmt.Errorf("TaskName is required")
	}

	if c.EnvironmentID == "" {
		return fmt.Errorf("EnvironmentID is required")
	}

	if c.DeployID == "" {
		return fmt.Errorf("DeployID is required")
	}

	return nil
}

func (c CreateTaskRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"task_name":      swagger.NewStringProperty(),
			"environment_id": swagger.NewStringProperty(),
			"deploy_id":      swagger.NewStringProperty(),
			// TODO: Change to Map defintion
			"container_overrides": swagger.NewStringProperty(),
		},
	}
}
