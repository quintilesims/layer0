package models

type Deploy struct {
	DeployFile []byte `json:"deploy_file"`
	DeployID   string `json:"deploy_id"`
	DeployName string `json:"deploy_name"`
	Version    string `json:"version"`
}
