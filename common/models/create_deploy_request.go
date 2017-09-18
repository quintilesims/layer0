package models

import (
	"fmt"
	"strings"
)

type CreateDeployRequest struct {
	DeployName string `json:"deploy_name"`
	DeployFile []byte `json:"deploy_file"`
}

func (c CreateDeployRequest) Validate() error {
	if c.DeployName == "" {
		return fmt.Errorf("Deploy Name is required")
	}

	if strings.Contains(c.DeployName, ".") {
		return fmt.Errorf("Deploy names cannot contain '.'")
	}

	if c.DeployFile == nil {
		return fmt.Errorf("Deploy file is required")
	}

	return nil
}
