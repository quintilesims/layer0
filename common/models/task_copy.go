package models

type TaskCopy struct {
	Details    []TaskDetail `json:"details"`
	Reason     string       `json:"reason"`
	TaskCopyID string       `json:"task_copy_id"`
}
