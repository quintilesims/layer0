package config

import "time"

const (
	DefaultAWSRegion           = "us-west-2"
	DefaultTimeBetweenRequests = time.Millisecond * 10
	DefaultJobExpiry           = time.Hour * 1
	DefaultLockExpiry          = time.Hour * 1
	DefaultPort                = 9090
	DefaultEnvironmentInstanceType     = "t2.small"
	DefaultEnvironmentMaxScale = 100
	  DefaultEnvironmentOS = "linux"
        DefaultServiceScale  = 1
)
