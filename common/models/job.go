package models

import (
	"time"

	"github.com/quintilesims/layer0/common/job"
)

type Job struct {
	JobID       string            `json:"job_id"`
	TaskID      string            `json:"task_id"`
	JobStatus   job.JobStatus     `json:"job_status"`
	JobType     job.JobType       `json:"job_type"`
	Request     string            `json:"request"`
	TimeCreated time.Time         `json:"time_created"`
	Meta        map[string]string `json:"meta"`
}
