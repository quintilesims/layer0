package models

type ServiceSummary struct {
	ServiceID       string `json:"service_id"`
	ServiceName     string `json:"service_name"`
	EnvironmentID   string `json:"environment_id"`
	EnvironmentName string `json:"environment_name"`
}
