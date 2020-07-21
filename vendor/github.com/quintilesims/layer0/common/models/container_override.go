package models

type ContainerOverride struct {
	ContainerName        string            `json:"container_name"`
	EnvironmentOverrides map[string]string `json:"environment_overrides"`
}
