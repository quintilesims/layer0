package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type Deploy struct {
	DeployFile []byte `json:"deploy_file"`
	DeployID   string `json:"deploy_id"`
	DeployName string `json:"deploy_name"`
	Version    string `json:"version"`
}

func (d Deploy) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"deploy_file": swagger.NewStringProperty(),
			"deploy_id":   swagger.NewStringProperty(),
			"deploy_name": swagger.NewStringProperty(),
			"version":     swagger.NewStringProperty(),
		},
	}
}
