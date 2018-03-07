package models

import swagger "github.com/zpatrick/go-plugin-swagger"

const (
	WindowsOS = "windows"
	LinuxOS   = "linux"
)

type Environment struct {
	EnvironmentID   string   `json:"environment_id"`
	EnvironmentName string   `json:"environment_name"`
	CurrentScale    int      `json:"current_size"`
	DesiredScale    int      `json:"desired_size"`
	InstanceType    string   `json:"instance_type"`
	SecurityGroupID string   `json:"security_group_id"`
	OperatingSystem string   `json:"operating_system"`
	AMIID           string   `json:"ami_id"`
	Links           []string `json:"links"`
}

func (e Environment) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"environment_id":    swagger.NewStringProperty(),
			"environment_name":  swagger.NewStringProperty(),
			"current_scale":     swagger.NewIntProperty(),
			"desired_scale":     swagger.NewIntProperty(),
			"instance_type":     swagger.NewStringProperty(),
			"security_group_id": swagger.NewStringProperty(),
			"operating_system":  swagger.NewStringProperty(),
			"ami_id":            swagger.NewStringProperty(),
			"links":             swagger.NewStringSliceProperty(),
		},
	}
}
