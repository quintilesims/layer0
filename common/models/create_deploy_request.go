package models

import "fmt"

type CreateDeployRequest struct {
	DeployName string `json:"deploy_name"`
	Dockerrun  []byte `json:"dockerrun"`
}

func (c CreateDeployRequest) Validate() error {
	if c.DeployName == "" {
		return fmt.Errorf("Deploy Name is required")
	}

	return nil
}
