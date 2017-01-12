package models

type TaskSummary struct {
	TaskID          string `json:"task_id"`
	TaskName        string `json:"task_name"`
	EnvironmentID   string `json:"environment_id"`
	EnvironmentName string `json:"environment_name"`
}
