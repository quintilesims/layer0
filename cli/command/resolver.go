package command

type Resolver interface {
	Resolve(entityType, target string) ([]string, error)
}
