package models

import (
	"time"

	swagger "github.com/zpatrick/go-plugin-swagger"
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
