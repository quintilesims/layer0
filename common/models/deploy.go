package models

type Deploy struct {
	DeployFile []byte `json:"dockerrun"`
	DeployID   string `json:"deploy_id"`
	DeployName string `json:"deploy_name"`
	Version    string `json:"version"`
}
