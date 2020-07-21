package models

type DeploySummary struct {
	DeployID   string `json:"deploy_id"`
	DeployName string `json:"deploy_name"`
	Version    string `json:"version"`
}
