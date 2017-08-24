package models

import (
	"time"
)

type Job struct {
	JobID       string            `json:"job_id"`
	TaskID      string            `json:"task_id"`
	JobStatus   string            `json:"job_status"`
	JobType     string            `json:"job_type"`
	Request     string            `json:"request"`
	TimeCreated time.Time         `json:"time_created"`
	Meta        map[string]string `json:"meta"`
}
