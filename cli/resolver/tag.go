package resolver

import (
	"github.com/quintilesims/layer0/client"
)

type TagResolver struct {
	client client.Client
}

func NewTagResolver(client client.Client) *TagResolver {
	return &TagResolver{
		client: client,
	}
}

func (r *TagResolver) Resolve(entityType, target string) ([]string, error) {
	// todo: implement
	return []string{target}, nil
}
