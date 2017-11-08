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

	// mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	// todo: put some old tags in the store, mock the ListTasks call
	tags := models.Tags{
		{
			EntityID:   "tsk_id1",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk_name1",
		},
		{
			EntityID:   "tsk_id1",
			EntityType: "task",
			Key:        "environment_id",
			Value:      "env_id1",
		},
		{
			EntityID:   "tsk_id1",
			EntityType: "task",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task/arn1",
		},
		{
			EntityID:   "tsk_id2",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk_name2",
		},
		{
			EntityID:   "tsk_id2",
			EntityType: "task",
			Key:        "environment_id",
			Value:      "env_id2",
		},
		{
			EntityID:   "tsk_id2",
			EntityType: "task",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task/arn2",
		},
		{
			EntityID:   "env_id1",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name1",
		},
		{
			EntityID:   "env_id2",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name2",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	janitor := NewJanitor(tagStore, taskProvider)
	if err := janitor.Run(); err != nil {
		t.Fatal(err)
	}

	// taskProvider.EXPECT().List().
	// 	Return(nil)
	// mockAWS.ECS.EXPECT().
	// 	ListClustersPages(&ecs.ListClustersInput{}, gomock.Any()).
	// 	Do(listClusterPagesFN).
	// 	Return(nil)

	expected := []models.TaskSummary{
		{
			TaskID:          "tsk_id1",
			TaskName:        "tsk_name1",
			EnvironmentID:   "env_id1",
			EnvironmentName: "env_name1",
		},
		{
			TaskID:          "tsk_id2",
			TaskName:        "tsk_name2",
			EnvironmentID:   "env_id2",
			EnvironmentName: "env_name2",
		},
	}
	actual := taskProvider.EXPECT().List().Return(expected, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	assert.Equal(t, expected, actual)
	// todo: assert only old tags got deleted
}
