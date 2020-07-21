package models

type CreateDeployRequest struct {
	DeployName string `json:"deploy_name"`
	Dockerrun  []byte `json:"dockerrun"`
}
