package models

import (
	"time"
)

type Job struct {
	JobID   string            `json:"job_id"`
	Type    string            `json:"type"`
	Request string            `json:"request"`
	Status  string            `json:"status"`
	Created time.Time         `json:"created"`
	Error   string            `json:"error"`
	Meta    map[string]string `json:"meta"`
}
