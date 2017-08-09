package models

import "github.com/quintilesims/layer0/common/job"

type CreateJobRequest struct {
	JobType job.JobType `json:"job_type"`
	Request string      `json:"request"`
}
