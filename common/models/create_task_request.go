package models

type CreateTaskRequest struct {
	ContainerOverrides []ContainerOverride `json:"container_overrides"`
	Copies             int64               `json:"copies"`
	DeployID           string              `json:"deploy_id"`
	EnvironmentID      string              `json:"environment_id"`
	TaskName           string              `json:"task_name"`
	DisableLogging     bool                `json:"disable_logging"`
}
