package models

type UpdateServiceRequest struct {
	DeployID       string `json:"deploy_id"`
	DisableLogging bool   `json:"disable_logging"`
}
