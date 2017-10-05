package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type APIConfig struct {
	Instance       string   `json:"instance"`
	VPCID          string   `json:"vpc_id"`
	Version        string   `json:"version"`
	PublicSubnets  []string `json:"public_subnets"`
	PrivateSubnets []string `json:"private_subnets"`
}

func (a APIConfig) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"instance":        swagger.NewStringProperty(),
			"vpc_id":          swagger.NewStringProperty(),
			"version":         swagger.NewStringProperty(),
			"public_subnets":  swagger.NewStringSliceProperty(),
			"private_subnets": swagger.NewStringSliceProperty(),
		},
	}
}
