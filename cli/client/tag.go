package client

import (
	"fmt"
	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) GetTags(params map[string]string) ([]*models.EntityWithTags, error) {
	query := "?"
	for k, v := range params {
		query += fmt.Sprintf("&%s=%s", k, v)
	}

	var response []*models.EntityWithTags
	if err := c.Execute(c.Sling("/tag").Get(query), &response); err != nil {
		return nil, err
	}

	return response, nil
}
