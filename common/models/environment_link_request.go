package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateEnvironmentLinkRequest EnvironmentLinkRequest

type DeleteEnvironmentLinkRequest EnvironmentLinkRequest

type EnvironmentLinkRequest struct {
	SourceEnvironmentID string `json:"source_environment_id"`
	DestEnvironmentID   string `json:"dest_environment_id"`
}

func (e EnvironmentLinkRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"source_environment_id": swagger.NewStringProperty(),
			"dest_environment_id":   swagger.NewStringProperty(),
		},
	}
}

func (e EnvironmentLinkRequest) Validate() error {
	if e.SourceEnvironmentID == "" {
		return fmt.Errorf("SourceEnvironmentID must be specified")
	}

	if e.DestEnvironmentID == "" {
		return fmt.Errorf("DestEnvironmentID must be specified")
	}

	return nil
}
