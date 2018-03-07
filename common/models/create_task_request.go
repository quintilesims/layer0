package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateTaskRequest struct {
	ContainerOverrides []ContainerOverride `json:"container_overrides"`
	DeployID           string              `json:"deploy_id"`
	EnvironmentID      string              `json:"environment_id"`
	TaskName           string              `json:"task_name"`
	Stateful           bool                `json:"stateful"`
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

	for _, o := range c.ContainerOverrides {
		if err := o.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c CreateTaskRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			// TODO: Change to Map defintion
			"container_overrides": swagger.NewStringProperty(),
			"deploy_id":           swagger.NewStringProperty(),
			"environment_id":      swagger.NewStringProperty(),
			"task_name":           swagger.NewStringProperty(),
			"stateful":            swagger.NewBoolProperty(),
		},
	}
}
