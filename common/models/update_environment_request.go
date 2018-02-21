package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type UpdateEnvironmentRequest struct {
	Scale *int      `json:"scale"`
	Links *[]string `json:"links"`
}

func (u UpdateEnvironmentRequest) Validate() error {
	if u.Scale != nil && *u.Scale < 0 {
		return fmt.Errorf("Scale must be omitted or a positive integer")
	}

	return nil
}

func (u UpdateEnvironmentRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"scale": swagger.NewIntProperty(),
			"links": swagger.NewStringSliceProperty(),
		},
	}
}
