package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentLink_createTags(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	environment := NewEnvironmentProvider(nil, tagStore, nil)

	if err := environment.createLinkTags("env_id1", "env_id2"); err != nil {
		t.Fatal(err)
	}

	expectedTags := models.Tags{
		{
			EntityID:   "env_id1",
			EntityType: "environment",
			Key:        "link",
			Value:      "env_id2",
		},
		{
			EntityID:   "env_id2",
			EntityType: "environment",
			Key:        "link",
			Value:      "env_id1",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}
