package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateEnvironmentRequest struct {
	EnvironmentName  string `json:"environment_name"`
	OperatingSystem  string `json:"operating_system"`
	Scale            int    `json:"scale"`
	InstanceType     string `json:"instance_type"`
	UserDataTemplate []byte `json:"user_data_template"`
	AMIID            string `json:"ami_id"`
}

func (r CreateEnvironmentRequest) Validate() error {
	if r.EnvironmentName == "" {
		return fmt.Errorf("EnvironmentName is required")
	}

	switch r.OperatingSystem {
	case LinuxOS:
	case "":
		return fmt.Errorf("OperatingSystem is required")
	default:
		return fmt.Errorf("Unrecognized OperatingSystem '%s'", r.OperatingSystem)
	}

	if r.Scale < 0 {
		return fmt.Errorf("Scale must be a positive integer")
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
			"scale":              swagger.NewIntProperty(),
			"operating_system":   swagger.NewStringProperty(),
			"ami_id":             swagger.NewStringProperty(),
		},
	}
}
