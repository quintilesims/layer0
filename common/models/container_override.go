package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type ContainerOverride struct {
	ContainerName        string            `json:"container_name"`
	EnvironmentOverrides map[string]string `json:"environment_overrides"`
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
