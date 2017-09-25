package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTask_deleteEntityTags(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	task := NewTaskProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "t1",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk1",
		},
		{
			EntityID:   "t1",
			EntityType: "task",
			Key:        "arn",
			Value:      "arn:aws:ecs:us-west-2:012345678910:task/aaaaaaaa-bbbb-cccc-eeee-ffffffffffff",
		},
		{
			EntityID:   "t2",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk2",
		},
		{
			EntityID:   "t2",
			EntityType: "task",
			Key:        "arn",
			Value:      "arn:aws:ecs:us-west-2:012345678910:task/00000000-1111-2222-3333-444444444444",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	for _, tag := range tags {
		if err := deleteEntityTags(task.TagStore, "task", tag.EntityID); err != nil {
			t.Fatal(err)
		}
	}

	assert.Len(t, tagStore.Tags(), 0)
}
