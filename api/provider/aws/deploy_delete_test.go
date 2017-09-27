package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDeploy_deleteDeployTags(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	deploy := NewDeployProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "deploy_id",
			EntityType: "deploy",
			Key:        "name",
			Value:      "deploy_name",
		},
		{
			EntityID:   "deploy_id",
			EntityType: "deploy",
			Key:        "version",
			Value:      "deploy_version",
		},
		{
			EntityID:   "deploy_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "deploy_arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	for _, tag := range tags {
		if err := deploy.deleteDeployTags(tag.EntityID); err != nil {
			t.Fatal(err)
		}
	}

	assert.Len(t, tagStore.Tags(), 0)
}
