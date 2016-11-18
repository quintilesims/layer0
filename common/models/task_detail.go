package models

type TaskDetail struct {
	ContainerName string `json:"container_name"`
	ExitCode      int64  `json:"exit_code"`
	LastStatus    string `json:"last_status"`
	Reason        string `json:"reason"`
}
