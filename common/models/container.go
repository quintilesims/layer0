package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type Container struct {
	ContainerName string `json:"container_name"`
	Status        string `json:"status"`
	ExitCode      int    `json:"exit_code"`
	Meta          string `json:"meta"`
}

func (t Container) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"container_name": swagger.NewStringProperty(),
			"status":         swagger.NewStringProperty(),
			"exit_code":      swagger.NewIntProperty(),
			"meta":           swagger.NewStringProperty(),
		},
	}
}
