package aws

import (
	"testing"

	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment_deleteTags(t *testing.T) {
	tagStore := tag_store.NewMemoryTagStore()
	environment := NewEnvironmentProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "eid1",
			EntityType: "environment",
			Key:        "name",
			Value:      "ename1",
		},
		{
			EntityID:   "eid1",
			EntityType: "environment",
			Key:        "os",
			Value:      "eos1",
		},
		{
			EntityID:   "eid2",
			EntityType: "environment",
			Key:        "name",
			Value:      "ename2",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	if err := environment.deleteTags("eid1"); err != nil {
		t.Fatal(err)
	}

	assert.NotContains(t, tagStore.Tags(), tags[0])
	assert.NotContains(t, tagStore.Tags(), tags[1])
	assert.Contains(t, tagStore.Tags(), tags[2])
}
