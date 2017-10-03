package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestService_deleteEntityTags(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	task := NewServiceProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "s1",
			EntityType: "service",
			Key:        "name",
			Value:      "svc1",
		},
		{
			EntityID:   "s2",
			EntityType: "service",
			Key:        "name",
			Value:      "svc2",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	for _, tag := range tags {
		if err := deleteEntityTags(task.TagStore, "service", tag.EntityID); err != nil {
			t.Fatal(err)
		}
	}

	assert.Len(t, tagStore.Tags(), 0)
}
