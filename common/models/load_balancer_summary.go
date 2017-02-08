package models

type LoadBalancerSummary struct {
	LoadBalancerID   string `json:"load_balancer_id"`
	LoadBalancerName string `json:"load_balancer_name"`
	EnvironmentID    string `json:"environment_id"`
	EnvironmentName  string `json:"environment_name"`
}
