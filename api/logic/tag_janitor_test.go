package logic

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestTagJanitorPulse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskLogicMock := mock_logic.NewMockTaskLogic(ctrl)
	tagStore, tagsAdded := getTagStore()

	tasks := []*models.TaskSummary{
		{
			TaskID: "task2",
		},
	}

	taskLogicMock.EXPECT().
		ListTasks().
		Return(tasks, nil)

	janitor := NewTagJanitor(taskLogicMock, tagStore)
	if err := janitor.pulse(); err != nil {
		t.Fatal(err)
	}

	tags, err := tagStore.SelectByType("task")
	if err != nil {
		t.Fatal(err)
	}

	//expecting 3 tags for task1 to be deleted from the tagstore
	testutils.AssertEqual(t, tagsAdded-3, len(tags))
	testutils.AssertEqual(t, "task2", tags[0].EntityID)
}

func getTagStore() (tag_store.TagStore, int) {
	store := tag_store.NewMemoryTagStore()
	tagsAdded := 0

	generateTag := func(entityType, entityID, key, value string) models.Tag {
		tagsAdded++

		return models.Tag{
			EntityType: entityType,
			EntityID:   entityID,
			Key:        key,
			Value:      value,
		}
	}

	//task 1
	store.Insert(generateTag("task", "task1", "deploy_id", "efgh-task.1"))
	store.Insert(generateTag("task", "task1", "environment_id", "env1"))
	store.Insert(generateTag("task", "task1", "name", "random-task"))

	//task 2
	store.Insert(generateTag("task", "task2", "deploy_id", "abcd-task.2"))
	store.Insert(generateTag("task", "task2", "environment_id", "env1"))
	store.Insert(generateTag("task", "task2", "name", "not-so-random-task"))

	return store, tagsAdded
}
