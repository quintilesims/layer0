package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTask_populateSummariesFromTaskARNs(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	task := NewTaskProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "t1",
			EntityType: "task",
			Key:        "name",
			Value:      "tname1",
		},
		{
			EntityID:   "t1",
			EntityType: "task",
			Key:        "arn",
			Value:      "arn1",
		},
		{
			EntityID:   "t1",
			EntityType: "task",
			Key:        "environment_id",
			Value:      "e1",
		},
		{
			EntityID:   "t2",
			EntityType: "task",
			Key:        "name",
			Value:      "tname2",
		},
		{
			EntityID:   "t2",
			EntityType: "task",
			Key:        "arn",
			Value:      "arn2",
		},
		{
			EntityID:   "t2",
			EntityType: "task",
			Key:        "environment_id",
			Value:      "e2",
		},
		{
			EntityID:   "e1",
			EntityType: "environment",
			Key:        "name",
			Value:      "ename1",
		},
		{
			EntityID:   "e2",
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

	arns := []string{"arn1", "arn2"}
	result, err := task.populateSummariesFromTaskARNs(arns)
	if err != nil {
		t.Fatal(err)
	}

	expected := []models.TaskSummary{
		{
			TaskID:          "t1",
			TaskName:        "tname1",
			EnvironmentID:   "e1",
			EnvironmentName: "ename1",
		},
		{
			TaskID:          "t2",
			TaskName:        "tname2",
			EnvironmentID:   "e2",
			EnvironmentName: "ename2",
		},
	}

	assert.Equal(t, expected, result)
}
