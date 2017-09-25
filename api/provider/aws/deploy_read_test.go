package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDeploy_populateModelTags(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	deploy := NewDeployProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "deploy_id_1",
			EntityType: "deploy",
			Key:        "name",
			Value:      "deploy_name_1",
		},
		{
			EntityID:   "deploy_id_1",
			EntityType: "deploy",
			Key:        "version",
			Value:      "deploy_version_1",
		},
		{
			EntityID:   "deploy_id_2",
			EntityType: "deploy",
			Key:        "name",
			Value:      "deploy_name_2",
		},
		{
			EntityID:   "deploy_id_2",
			EntityType: "deploy",
			Key:        "version",
			Value:      "deploy_version_2",
		},
		{
			EntityID:   "junk_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "junk_value",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	result := &models.Deploy{}
	if err := deploy.populateModelTags("deploy_id_1", result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "deploy_name_1", result.DeployName)
	assert.Equal(t, "deploy_version_1", result.Version)
}
