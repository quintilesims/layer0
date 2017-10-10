package resolver

import (
	"github.com/quintilesims/layer0/client"
)

func NewTagResolver(client client.Client) ResolverFunc {
	return func(entityType, target string) ([]string, error) {
		return []string{target}, nil
	}
}
