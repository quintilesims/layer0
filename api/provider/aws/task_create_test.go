package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTask_createTags(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	task := NewTaskProvider(nil, tagStore, nil)

	if err := task.createTags("tsk_id", "tsk_name", "env_id", "tsk_arn"); err != nil {
		t.Fatal(err)
	}

	expectedTags := models.Tags{
		{
			EntityID:   "tsk_id",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk_name",
		},
		{
			EntityID:   "tsk_id",
			EntityType: "task",
			Key:        "environment_id",
			Value:      "env_id",
		},
		{
			EntityID:   "tsk_id",
			EntityType: "task",
			Key:        "arn",
			Value:      "tsk_arn",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}
