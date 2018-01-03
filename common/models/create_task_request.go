package models

type CreateTaskRequest struct {
	ContainerOverrides []ContainerOverride `json:"container_overrides"`
	DeployID           string              `json:"deploy_id"`
	EnvironmentID      string              `json:"environment_id"`
	TaskName           string              `json:"task_name"`
}
