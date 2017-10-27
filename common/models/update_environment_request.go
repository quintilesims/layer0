package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type UpdateEnvironmentRequestJob struct {
	EnvironmentID string `json:"environment_id"`
	UpdateEnvironmentRequest
}

type UpdateEnvironmentRequest struct {
	MinClusterCount *int      `json:"min_cluster_count"`
	Links           *[]string `json:"links"`
}

func (u UpdateEnvironmentRequest) Validate() error {
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
			"min_cluster_count": swagger.NewIntProperty(),
			"links":             swagger.NewStringSliceProperty(),
		},
	}
}
