package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type ContainerOverride struct {
	ContainerName        string            `json:"container_name"`
	EnvironmentOverrides map[string]string `json:"environment_overrides"`
}

func (c ContainerOverride) Validate() error {
	if c.ContainerName == "" {
		return fmt.Errorf("ContainerName is required")
	}

	return nil
}

func (c ContainerOverride) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"container_name":        swagger.NewStringProperty(),
			"environment_overrides": swagger.NewStringSliceProperty(),
		},
	}
}
