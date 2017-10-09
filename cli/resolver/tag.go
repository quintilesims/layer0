package resolver

import (
	"fmt"

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
	return nil, fmt.Errorf("TagResolver.Resolve not implemented!")
}
