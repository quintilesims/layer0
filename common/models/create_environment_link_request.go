package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type CreateEnvironmentLinkRequest struct {
	EnvironmentID string `json:"environment_id"`
}

func (c CreateEnvironmentLinkRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"environment_id": swagger.NewStringProperty(),
		},
	}
}
