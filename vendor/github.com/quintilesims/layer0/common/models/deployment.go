package models

import (
	"time"
)

type Deployment struct {
	Created       time.Time `json:"created"`
	DeployID      string    `json:"deploy_id"`
	DeployName    string    `json:"deploy_name"`
	DeployVersion string    `json:"deploy_version"`
	DesiredCount  int64     `json:"desired_count"`
	PendingCount  int64     `json:"pending_count"`
	RunningCount  int64     `json:"running_count"`
	Status        string    `json:"status"`
	Updated       time.Time `json:"updated"`
	DeploymentID  string    `json:"deployment_id"`
}
