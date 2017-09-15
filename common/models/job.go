package models

import (
	"time"
)

type Job struct {
	JobID   string    `json:"job_id"`
	Type    string    `json:"type"`
	Status  string    `json:"status"`
	Request string    `json:"request"`
	Result  string    `json:"result"`
	Created time.Time `json:"created"`
	Error   string    `json:"error"`
}
