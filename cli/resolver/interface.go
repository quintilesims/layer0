package resolver

type Resolver interface {
	Resolve(entityType, target string) ([]string, error)
}

type ResolverFunc func(entityType, target string) ([]string, error)

func (r ResolverFunc) Resolve(entityType, target string) ([]string, error) {
	return r(entityType, target)
}
