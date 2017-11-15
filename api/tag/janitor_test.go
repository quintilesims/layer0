package tag

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestJanitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagStore := NewMemoryStore()
	taskProvider := mock_provider.NewMockTaskProvider(ctrl)

	tags := models.Tags{
		{
			EntityID:   "tsk_id1",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk_name1",
		},
		{
			EntityID:   "tsk_id2",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk_name2",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	tasks := []models.TaskSummary{
		{
			TaskID: "tsk_id2",
		},
	}

	taskProvider.EXPECT().
		List().
		Return(tasks, nil)

	janitor := NewJanitor(tagStore, taskProvider)
	if err := janitor.Run(); err != nil {
		t.Fatal(err)
	}

	actual, err := tagStore.SelectByType("task")
	if err != nil {
		t.Fatal(err)
	}

	expected := models.Tags{
		{
			EntityID:   "tsk_id2",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk_name2",
		},
	}

	assert.Equal(t, expected, actual)
}
