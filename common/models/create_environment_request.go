package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateEnvironmentRequest struct {
	EnvironmentName  string `json:"environment_name"`
	InstanceSize     string `json:"instance_size"`
	UserDataTemplate []byte `json:"user_data_template"`
	MinClusterCount  int    `json:"min_cluster_count"`
	OperatingSystem  string `json:"operating_system"`
	AMIID            string `json:"ami_id"`
}

func (r CreateEnvironmentRequest) Validate() error {
	if r.EnvironmentName == "" {
		return fmt.Errorf("EnvironmentName is required")
	}

	if r.MinClusterCount < 0 {
		return fmt.Errorf("MinClusterCount must be a positive integer")
	}

	if r.OperatingSystem == "" {
		return fmt.Errorf("OperatingSystem is required")
	}

	return nil
}

func (e CreateEnvironmentRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"environment_name":   swagger.NewStringProperty(),
			"instance_size":      swagger.NewStringProperty(),
			"user_data_template": swagger.NewStringProperty(),
			"min_cluster_count":  swagger.NewIntProperty(),
			"operating_system":   swagger.NewStringProperty(),
			"ami_id":             swagger.NewStringProperty(),
		},
	}
}
