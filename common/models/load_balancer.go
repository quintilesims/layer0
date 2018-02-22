package models

type LoadBalancer struct {
	EnvironmentID    string      `json:"environment_id"`
	EnvironmentName  string      `json:"environment_name"`
	HealthCheck      HealthCheck `json:"health_check"`
	IdleTimeout      int64       `json:"idle_timeout"`
	IsPublic         bool        `json:"is_public"`
	LoadBalancerID   string      `json:"load_balancer_id"`
	LoadBalancerName string      `json:"load_balancer_name"`
	Ports            []Port      `json:"ports"`
	ServiceID        string      `json:"service_id"`
	ServiceName      string      `json:"service_name"`
	URL              string      `json:"url"`
}
