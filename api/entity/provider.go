package entity

type Provider interface {
	ListEnvironmentIDs() ([]string, error)
	GetEnvironment(environmentID string) Environment

	ListJobIDs() ([]string, error)
	GetJob(jobID string) Job
}
