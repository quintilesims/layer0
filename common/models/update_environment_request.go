package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type UpdateEnvironmentRequest struct {
	EnvironmentID   string `json:"environment_id"`
	MinClusterCount *int   `json:"min_cluster_count"`
}

func (u UpdateEnvironmentRequest) Validate() error {
	if u.EnvironmentID == "" {
		return fmt.Errorf("EnvironmentID must be specified")
	}

	if u.MinClusterCount != nil {
		if *u.MinClusterCount < 0 {
			return fmt.Errorf("MinClusterCount must be a positive integer")
		}
	}

	return nil
}

func (u UpdateEnvironmentRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"environment_id":    swagger.NewStringProperty(),
			"min_cluster_count": swagger.NewIntProperty(),
		},
	}
}
