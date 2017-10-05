package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type APIConfig struct {
	Prefix         string   `json:"prefix"`
	VPCID          string   `json:"vpc_id"`
	Version        string   `json:"version"`
	PublicSubnets  []string `json:"public_subnets"`
	PrivateSubnets []string `json:"private_subnets"`
}

func (a APIConfig) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"prefix":          swagger.NewStringProperty(),
			"vpc_id":          swagger.NewStringProperty(),
			"public_subnets":  swagger.NewStringSliceProperty(),
			"private_subnets": swagger.NewStringSliceProperty(),
		},
	}
}
