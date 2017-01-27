package models

type UpdateLoadBalancerPortsRequest struct {
    Ports []Port `json:"ports"`
}
