package types

import (
	"fmt"
)

type JobStatus int64

const (
	Pending JobStatus = 1 + iota
	InProgress
	Completed
	Error
)

var jobStatusStrings = []string{
	"pending",
	"in progress",
	"completed",
	"error",
}

func (jobStatus JobStatus) String() string {
	if int(jobStatus) < 1 || int(jobStatus) > len(jobStatusStrings) {
		return fmt.Sprintf("(Unknown JobStatus=%d!)", jobStatus)
	}

	return jobStatusStrings[jobStatus-1]
}

type JobType int64

const (
	DeleteEnvironmentJob JobType = 1 + iota
	DeleteServiceJob
	DeleteLoadBalancerJob
	DeleteTaskJob
	CreateTaskJob
)

var jobTypeStrings = []string{
	"delete environment",
	"delete service",
	"delete load balancer",
	"delete task",
	"create task",
}

func (jobType JobType) String() string {
	if int(jobType) < 1 || int(jobType) > len(jobTypeStrings) {
		return fmt.Sprintf("(Unknown JobType=%d!)", int(jobType))
	}

	return jobTypeStrings[jobType-1]
}
