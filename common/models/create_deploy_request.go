package models

import (
	"fmt"
	swagger "github.com/zpatrick/go-plugin-swagger"
)

type CreateDeployRequest struct {
	DeployName string `json:"deploy_name"`
	DeployFile []byte `json:"deploy_file"`
}

func (c CreateDeployRequest) Validate() error {
	if c.DeployName == "" {
		return fmt.Errorf("Deploy Name is required")
	}

	if len(c.DeployFile) == 0 {
		return fmt.Errorf("Deploy file is required")
	}

	return nil
}

func (c CreateDeployRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"deploy_name": swagger.NewStringProperty(),
			// TODO: it doesn't appear that byte slices exist?
			// "deploy_file": swagger.NewByteSliceProperty(),
		},
	}
}
