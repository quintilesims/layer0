package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateEnvironmentRequest struct {
	EnvironmentName  string `json:"environment_name"`
	EnvironmentType  string `json:"environment_type"`
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

	switch r.EnvironmentType {
	case EnvironmentTypeStatic, EnvironmentTypeDynamic:
	case "":
		return fmt.Errorf("EnvironmentType is required")
	default:
		return fmt.Errorf("Unrecognized environment type '%s'", r.EnvironmentType)
	}

	switch r.OperatingSystem {
	case LinuxOS, WindowsOS:
	case "":
		return fmt.Errorf("OperatingSystem is required")
	default:
		return fmt.Errorf("Unrecognized OperatingSystem '%s'", r.OperatingSystem)
	}

	if r.EnvironmentType == EnvironmentTypeDynamic {
		if r.OperatingSystem != LinuxOS {
			return fmt.Errorf("Only the '%s' OperatingSystem can be specified with dynamic environments", LinuxOS)
		}

		if r.Scale != 0 {
			return fmt.Errorf("Cannot specify scale with dynamic environments")
		}

		if r.InstanceType != "" {
			return fmt.Errorf("Cannot specify InstanceType with dynamic environments")
		}

		if len(r.UserDataTemplate) != 0 {
			return fmt.Errorf("Cannot specify UserDataTemplate with dynamic environments")
		}

		if r.AMIID != "" {
			return fmt.Errorf("Cannot specify AMI ID with dynamic environments")
		}
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
			"environment_type":   swagger.NewStringProperty(),
			"user_data_template": swagger.NewStringProperty(),
			"scale":              swagger.NewIntProperty(),
			"operating_system":   swagger.NewStringProperty(),
			"ami_id":             swagger.NewStringProperty(),
		},
	}
}
