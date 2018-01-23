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

	if strings.EqualFold(r.EnvironmentType, EnvironmentTypeDynamic) &&
		!strings.EqualFold(r.OperatingSystem, LinuxOS) {
		return fmt.Errorf("%s is not a supported OS for dynamic environments", r.OperatingSystem)
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
