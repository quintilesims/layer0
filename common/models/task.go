package models

type Task struct {
	Copies          []TaskCopy `json:"copies"`
	DeployID        string     `json:"deploy_id"`
	DeployName      string     `json:"deploy_name"`
	DeployVersion   string     `json:"deploy_version"`
	DesiredCount    int64      `json:"desired_count"`
	EnvironmentID   string     `json:"environment_id"`
	EnvironmentName string     `json:"environment_name"`
	PendingCount    int64      `json:"pending_count"`
	RunningCount    int64      `json:"running_count"`
	TaskID          string     `json:"task_id"`
	TaskName        string     `json:"task_name"`
}
