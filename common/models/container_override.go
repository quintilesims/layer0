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

	// emtpy values are ok
	for k, _ := range c.EnvironmentOverrides {
		if k == "" {
			return fmt.Errorf("EnvironmentOverride is missing a key")
		}
	}

	return nil
}

func (t ContainerOverride) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"container_name": swagger.NewStringProperty(),
			//TODO: Change to Map
			"environment_overrides": swagger.NewStringProperty(),
		},
	}
}
