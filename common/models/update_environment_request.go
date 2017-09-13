package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type UpdateEnvironmentRequest struct {
	EnvironmentID   string `json:"environment_id"`
	MinClusterCount int64  `json:"min_cluster_count"`
}

func (r UpdateEnvironmentRequest) Validate() error {
	if r.EnvironmentID == "" {
		return fmt.Errorf("EnvironmentID must be specified")
	}

	if r.MinClusterCount < 0 {
		return fmt.Errorf("MinClusterCount must be a positive integer")
	}

	return nil
}

func (e UpdateEnvironmentRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"environment_id":    swagger.NewStringProperty(),
			"min_cluster_count": swagger.NewIntProperty(),
		},
	}
}
