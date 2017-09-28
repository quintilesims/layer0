package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type LogFile struct {
	ContainerName string   `json:"container_name"`
	Lines         []string `json:"lines"`
}

func (l LogFile) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"container_name": swagger.NewStringProperty(),
			"lines":          swagger.NewStringSliceProperty(),
		},
	}
}
