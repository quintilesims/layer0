package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type UpdateEnvironmentRequest struct {
	MinScale *int      `json:"min_scale"`
	MaxScale *int      `json:"max_scale"`
	Links    *[]string `json:"links"`
}

func (u UpdateEnvironmentRequest) Validate() error {
	if u.MinScale != nil && *u.MinScale < 0 {
		return fmt.Errorf("MinScale must be a positive integer")
	}

	if u.MaxScale != nil && *u.MaxScale < 0 {
		return fmt.Errorf("MaxScale must be a positive integer")
	}

	if u.MaxScale != nil && u.MinScale != nil && *u.MaxScale < *u.MinScale {
		return fmt.Errorf("MaxScale must be greater than or equal to MinScale")
	}

	return nil
}

func (u UpdateEnvironmentRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"min_scale": swagger.NewIntProperty(),
			"max_scale": swagger.NewIntProperty(),
			"links":     swagger.NewStringSliceProperty(),
		},
	}
}
