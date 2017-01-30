package models

type UpdateLoadBalancerHealthCheckRequest struct {
	HealthCheck HealthCheck `json:"health_check"`
}
