package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment_makeEnvironmentModels(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	environment := NewEnvironmentProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name",
		},
		{
			EntityID:   "env_id_junk",
			EntityType: "environment",
			Key:        "name",
			Value:      "invalid",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "os",
			Value:      "env_os",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "link",
			Value:      "env2",
		},
		{
			EntityID:   "env2",
			EntityType: "environment",
			Key:        "link",
			Value:      "env_id",
		},
		{
			EntityID:   "env_id",
			EntityType: "service",
			Key:        "os",
			Value:      "bad_os",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	result, err := environment.makeEnvironmentModel("env_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "env_name", result.EnvironmentName)
	assert.Equal(t, "env_os", result.OperatingSystem)
	assert.Equal(t, []string{"env2"}, result.Links)
}
