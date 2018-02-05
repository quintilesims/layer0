package models

import (
	"fmt"
	"strings"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateEnvironmentRequest struct {
	EnvironmentName  string `json:"environment_name"`
	EnvironmentType  string `json:"environment_type"`
	Scale            int    `json:"scale"`
	InstanceType     string `json:"instance_type"`
	UserDataTemplate []byte `json:"user_data_template"`
	OperatingSystem  string `json:"operating_system"`
	AMIID            string `json:"ami_id"`
}

func (r CreateEnvironmentRequest) Validate() error {
	if r.EnvironmentName == "" {
		return fmt.Errorf("EnvironmentName is required")
	}

	if r.Scale < 0 {
		return fmt.Errorf("Scale must be a positive integer")
	}

	if !strings.EqualFold(r.EnvironmentType, EnvironmentTypeDynamic) &&
		!strings.EqualFold(r.EnvironmentType, EnvironmentTypeStatic) {
		return fmt.Errorf("%s is not a supported/valid environment type", r.EnvironmentType)
	}

	if strings.EqualFold(r.EnvironmentType, EnvironmentTypeDynamic) &&
		!strings.EqualFold(r.OperatingSystem, LinuxOS) {
		return fmt.Errorf("%s is not a supported OS for dynamic environments", r.OperatingSystem)
	}

	if strings.EqualFold(r.EnvironmentType, EnvironmentTypeDynamic) &&
		r.Scale > 0 {
		return fmt.Errorf("setting `Scale` is not valid for dynamic environments")
	}

	if r.EnvironmentType == "" {
		return fmt.Errorf("EnvironmentType is required, please specify 'static' or 'dynamic'")
	}

	if r.InstanceType == "" {
		return fmt.Errorf("Instancetype is required see - https://aws.amazon.com/ec2/instance-types/")
	}

	if r.OperatingSystem == "" {
		return fmt.Errorf("OperatingSystem is required, please specify 'windows' or 'linux'")
	}

	if r.AMIID == "" {
		return fmt.Errorf("AMIID is required")
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
