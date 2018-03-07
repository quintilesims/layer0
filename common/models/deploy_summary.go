package models

type DeploySummary struct {
	Compatibilities []string `json:"compatibilities"`
	DeployID        string   `json:"deploy_id"`
	DeployName      string   `json:"deploy_name"`
	Version         string   `json:"version"`
}
