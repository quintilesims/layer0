package models

type UpdateLoadBalancerIdleTimeoutRequest struct {
	IdleTimeout int `json:"health_check"`
}
