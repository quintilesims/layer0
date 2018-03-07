package config

import (
	"time"

	"github.com/quintilesims/layer0/common/models"
)

const (
	DefaultAWSRegion               = "us-west-2"
	DefaultTimeBetweenRequests     = time.Millisecond * 10
	DefaultLockExpiry              = time.Hour * 1
	DefaultPort                    = 9090
	DefaultEnvironmentInstanceType = "t2.small"
	DefaultEnvironmentType         = models.EnvironmentTypeDynamic
	DefaultEnvironmentOS           = "linux"
	DefaultServiceScale            = 1
	DefaultRetryTimeOut            = time.Second * 60
	DefaultRetryWaitTime           = time.Second
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
