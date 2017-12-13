package config

import (
	"time"

	"github.com/quintilesims/layer0/common/models"
)

const (
	DefaultAWSRegion               = "us-west-2"
	DefaultTimeBetweenRequests     = time.Millisecond * 10
	DefaultJobExpiry               = time.Hour * 1
	DefaultLockExpiry              = time.Hour * 1
	DefaultPort                    = 9090
	DefaultEnvironmentInstanceType = "t2.small"
	DefaultEnvironmentMaxScale     = 100
	DefaultEnvironmentOS           = "linux"
	DefaultServiceScale            = 1
)

var DefaultLoadBalancerHealthCheck = models.HealthCheck{
	Target:             "TCP:80",
	Interval:           30,
	Timeout:            5,
	HealthyThreshold:   2,
	UnhealthyThreshold: 2,
}

var DefaultLoadBalancerPort = models.Port{
	ContainerPort: 80,
	HostPort:      80,
	Protocol:      "TCP",
}
