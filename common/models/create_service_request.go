package models

type CreateServiceRequest struct {
	DeployID       string `json:"deploy_id"`
	EnvironmentID  string `json:"environment_id"`
	LoadBalancerID string `json:"load_balancer_id"`
	ServiceName    string `json:"service_name"`
	DisableLogging bool   `json:"disable_logging"`
}
