package models

type UpdateLoadBalancerRequest struct {
	LoadBalancerID string       `json:"load_balancer_id"`
	HealthCheck    *HealthCheck `json:"health_check"`
	Ports          *[]Port      `json:"ports"`
}
