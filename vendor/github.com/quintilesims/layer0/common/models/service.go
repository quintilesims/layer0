package models

type Service struct {
	Deployments      []Deployment `json:"deployments"`
	DesiredCount     int64        `json:"desired_count"`
	EnvironmentID    string       `json:"environment_id"`
	EnvironmentName  string       `json:"environment_name"`
	LoadBalancerID   string       `json:"load_balancer_id"`
	LoadBalancerName string       `json:"load_balancer_name"`
	PendingCount     int64        `json:"pending_count"`
	RunningCount     int64        `json:"running_count"`
	ServiceID        string       `json:"service_id"`
	ServiceName      string       `json:"service_name"`
}
