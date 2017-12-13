package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateEnvironmentRequest struct {
	EnvironmentName  string `json:"environment_name"`
	InstanceType     string `json:"instance_type"`
	UserDataTemplate []byte `json:"user_data_template"`
	MinScale         int    `json:"min_scale"`
	MaxScale         int    `json:"max_scale"`
	OperatingSystem  string `json:"operating_system"`
	AMIID            string `json:"ami_id"`
}

func (r CreateEnvironmentRequest) Validate() error {
	if r.EnvironmentName == "" {
		return fmt.Errorf("EnvironmentName is required")
	}

	if r.MinScale < 0 {
		return fmt.Errorf("MinScale must be a positive integer")
	}

	if r.MaxScale < 0 {
		return fmt.Errorf("MaxScale must be a positive integer")
	}

	if r.MaxScale < r.MinScale {
		return fmt.Errorf("MaxScale must be greater than or equal to MinScale")
	}

	return nil
}

func (e CreateEnvironmentRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"environment_name":   swagger.NewStringProperty(),
			"instance_type":      swagger.NewStringProperty(),
			"user_data_template": swagger.NewStringProperty(),
			"min_scale":          swagger.NewIntProperty(),
			"max_scale":          swagger.NewIntProperty(),
			"operating_system":   swagger.NewStringProperty(),
			"ami_id":             swagger.NewStringProperty(),
		},
	}
}
