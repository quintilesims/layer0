package models

type EnvironmentSummary struct {
	EnvironmentID   string `json:"environment_id"`
	EnvironmentName string `json:"environment_name"`
	OperatingSystem string `json:"operating_system"`
}
