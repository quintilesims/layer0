package models

import (
	"time"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type Job struct {
	JobID   string    `json:"job_id"`
	Type    JobType   `json:"type"`
	Status  JobStatus `json:"status"`
	Request string    `json:"request"`
	Result  string    `json:"result"`
	Created time.Time `json:"created"`
	Error   string    `json:"error"`
}

func (j Job) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"job_id":  swagger.NewStringProperty(),
			"type":    swagger.NewStringProperty(),
			"status":  swagger.NewStringProperty(),
			"request": swagger.NewStringProperty(),
			"result":  swagger.NewStringProperty(),
			"created": swagger.NewStringProperty(),
			"error":   swagger.NewStringProperty(),
		},
	}
}

type JobType string

func (j JobType) String() string {
	return string(j)
}

const (
	CreateDeployJob       JobType = "CreateDeploy"
	CreateEnvironmentJob  JobType = "CreateEnvironment"
	LinkEnvironmentJob    JobType = "LinkEnvironment"
	UnlinkEnvironmentJob  JobType = "UnlinkEnvironment"
	CreateLoadBalancerJob JobType = "CreateLoadBalancer"
	CreateServiceJob      JobType = "CreateService"
	CreateTaskJob         JobType = "CreateTask"
	DeleteDeployJob       JobType = "DeleteDeploy"
	DeleteEnvironmentJob  JobType = "DeleteEnvironment"
	DeleteLoadBalancerJob JobType = "DeleteLoadBalancer"
	DeleteServiceJob      JobType = "DeleteService"
	DeleteTaskJob         JobType = "DeleteTask"
	UpdateEnvironmentJob  JobType = "UpdateEnvironment"
	UpdateLoadBalancerJob JobType = "UpdateLoadBalancer"
	UpdateServiceJob      JobType = "UpdateService"
)

type JobStatus string

func (j JobStatus) String() string {
	return string(j)
}

const (
	PendingJobStatus    JobStatus = "PendingJobStatus"
	InProgressJobStatus JobStatus = "InProgress"
	CompletedJobStatus  JobStatus = "Completed"
	ErrorJobStatus      JobStatus = "Error"
)
