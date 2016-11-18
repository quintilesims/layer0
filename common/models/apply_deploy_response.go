package models

type ApplyDeployResponse struct {
	DeployID  string `json:"deploy_id"`
	ServiceID string `json:"service_id"`
}
