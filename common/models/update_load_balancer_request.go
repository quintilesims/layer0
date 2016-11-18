package models

type UpdateLoadBalancerRequest struct {
	Ports []Port `json:"ports"`
}
