package models

import (
	"fmt"

	swagger "github.com/zpatrick/go-plugin-swagger"
)

type UpdateServiceRequest struct {
	ServiceID string  `json:"service_id"`
	DeployID  *string `json:"deploy_id"`
	Scale     *int    `json:"service_scale_count"`
}

func (u UpdateServiceRequest) Validate() error {
	if u.ServiceID == "" {
		return fmt.Errorf("ServiceID must be specified")
	}

	if u.DeployID != nil && *u.DeployID == "" {
		return fmt.Errorf("DeployID must be specified")
	}

	if u.Scale != nil && *u.Scale < 0 {
		return fmt.Errorf("Scale must be a positive integer")
	}

	return nil
}

func (u UpdateServiceRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"service_id": swagger.NewStringProperty(),
			"deploy_id":  swagger.NewStringProperty(),
			"scale":      swagger.NewIntProperty(),
		},
	}
}
