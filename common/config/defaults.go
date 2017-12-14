package config

import (
	"time"

	"github.com/quintilesims/layer0/common/models"
)

const (
	DefaultEndpoint                = "http://localhost:9090/"
	DefaultPort                    = 9090
	DefaultNumWorkers              = 10
	DefaultJobExpiry               = time.Hour
	DefaultLockExpiry              = time.Hour
	DefaultAWSRegion               = "us-west-2"
	DefaultAWSRequestDelay         = time.Millisecond
	DefaultOutput                  = "text"
	DefaultTimeout                 = time.Minute * 15
	DefaultEnvironmentInstanceType = "t2.small"
	DefaultEnvironmentMaxScale     = 100
	DefaultEnvironmentOS           = "linux"
	DefaultServiceScale            = 1
)

func DefaultLoadBalancerHealthCheck() models.HealthCheck {
	return models.HealthCheck{
		Target:             "TCP:80",
		Interval:           30,
		Timeout:            5,
		HealthyThreshold:   2,
		UnhealthyThreshold: 2,
	}
}

func DefaultLoadBalancerPort() models.Port {
	return models.Port{
		ContainerPort: 80,
		HostPort:      80,
		Protocol:      "TCP",
	}
}
