package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type UpdateServiceRequest struct {
	DeployID *string `json:"deploy_id"`
	Scale    *int    `json:"scale"`
}

func (u UpdateServiceRequest) Validate() error {
	if u.DeployID != nil && *u.DeployID == "" {
		return fmt.Errorf("DeployID must be omitted or non-empty string")
	}

	if u.Scale != nil && *u.Scale < 0 {
		return fmt.Errorf("Scale must be omitted or a positive integer")
	}

	return nil
}

func (u UpdateServiceRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"deploy_id": swagger.NewStringProperty(),
			"scale":     swagger.NewIntProperty(),
		},
	}
}
