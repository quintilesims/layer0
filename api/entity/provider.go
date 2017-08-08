package entity

type Provider interface {
	ListEnvironmentIDs() ([]string, error)
	GetEnvironment(environmentID string) Environment
}
