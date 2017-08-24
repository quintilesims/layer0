package models

type ScheduleJobRequest struct {
	JobType string      `json:"job_type"`
	Request interface{} `json:"request"`
}
