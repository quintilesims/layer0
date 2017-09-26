package models

import (
	"fmt"
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
