package models

import "fmt"

type CreateTaskRequest struct {
	ContainerOverrides []ContainerOverride `json:"container_overrides"`
	Copies             int                 `json:"copies"`
	DeployID           string              `json:"deploy_id"`
	EnvironmentID      string              `json:"environment_id"`
	TaskName           string              `json:"task_name"`
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
