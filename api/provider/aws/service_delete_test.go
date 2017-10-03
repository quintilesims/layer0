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

	tag := models.Tag{
		EntityID:   "s1",
		EntityType: "service",
		Key:        "name",
		Value:      "svc1",
	}

	if err := tagStore.Insert(tag); err != nil {
		t.Fatal(err)
	}

	if err := deleteEntityTags(task.TagStore, "service", tag.EntityID); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, tagStore.Tags(), 0)
}
