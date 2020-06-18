package models

type APIConfig struct {
	Prefix         string   `json:"prefix"`
	VPCID          string   `json:"vpc_id"`
	PublicSubnets  []string `json:"public_subnets"`
	PrivateSubnets []string `json:"private_subnets"`
}
