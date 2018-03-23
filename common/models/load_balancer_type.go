package models

import (
	"strings"
)

type LoadBalancerType string

const (
	ClassicLoadBalancerType     LoadBalancerType = "CLB"
	ApplicationLoadBalancerType LoadBalancerType = "ALB"
)

func validLoadBalancerTypes() map[LoadBalancerType]string {
	return map[LoadBalancerType]string{
		ClassicLoadBalancerType:     "classic",
		ApplicationLoadBalancerType: "application",
	}
}

func (t LoadBalancerType) String() string {
	key := strings.ToUpper(string(t))
	if v, ok := validLoadBalancerTypes()[LoadBalancerType(key)]; ok {
		return v
	}

	return "unknown"
}

func (t LoadBalancerType) Equals(s LoadBalancerType) bool {
	return strings.ToUpper(string(t)) == strings.ToUpper(string(s))
}

func (t LoadBalancerType) IsValid() bool {
	if _, ok := validLoadBalancerTypes()[t]; ok {
		return true
	}

	return false
}
