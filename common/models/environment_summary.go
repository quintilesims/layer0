package models

type EnvironmentSummary struct {
	EnvironmentID   string `json:"environment_id"`
	EnvironmentName string `json:"environment_name"`
	EnvironmentType string `json:"environment_type"`
	OperatingSystem string `json:"operating_system"`
}
