package client

import (
	"net/url"

	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/rclient"
)

func (c *APIClient) ListTags(query url.Values) (models.Tags, error) {
	var tags models.Tags
	if err := c.client.Get("/tags", &tags, rclient.Query(query)); err != nil {
		return nil, err
	}

	return tags, nil
}
