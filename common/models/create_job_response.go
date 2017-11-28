package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type CreateJobResponse struct {
	JobID string `json:"job_id"`
}

func (c CreateJobResponse) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"job_id": swagger.NewStringProperty(),
		},
	}
}
