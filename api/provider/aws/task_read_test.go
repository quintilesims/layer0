package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTask_makeTaskModel(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	task := NewTaskProvider(nil, tagStore, nil)

	tags := models.Tags{
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
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "name",
			Value:      "dpl_name",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "version",
			Value:      "dpl_version",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	result, err := task.makeTaskModel("tsk_id", "env_id", "dpl_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "tsk_name", result.TaskName)
	assert.Equal(t, "env_name", result.EnvironmentName)
	assert.Equal(t, "dpl_name", result.DeployName)
	assert.Equal(t, "dpl_version", result.DeployVersion)
}
