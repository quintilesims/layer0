package models

type ScalerRunInfo struct {
	EnvironmentID           string     `json:"environment_id"`
	ScaleBeforeRun          int        `json:"scale_before_run"`
	DesiredScaleAfterRun    int        `json:"desired_scale_after_run"`
	ActualScaleAfterRun     int        `json:"actual_scale_after_run"`
	UnusedResourceProviders int        `json:"unused_resource_providers"`
	PendingResources        []Resource `json:"pending_resources"`
	ResourceProviders       []Resource `json:"resource_providers"`
}
