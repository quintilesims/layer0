package logic

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

var taskTags = []models.Tag{
	{EntityID: "t1", EntityType: "task", Key: "name", Value: "task1"},
	{EntityID: "t1", EntityType: "task", Key: "environment_id", Value: "e1"},
	{EntityID: "t1", EntityType: "task", Key: "deploy_id", Value: "abcd.1"},
	{EntityID: "t2", EntityType: "task", Key: "name", Value: "task2"},
	{EntityID: "t2", EntityType: "task", Key: "environment_id", Value: "e2"},
	{EntityID: "t2", EntityType: "task", Key: "deploy_id", Value: "efgh.1"},
}

func TestTagJanitorPulse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskLogicMock := mock_logic.NewMockTaskLogic(ctrl)
	tagStore, tagsAdded := getTagStore(taskTags)

	tasks := []*models.TaskSummary{
		{
			TaskID: "t2",
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
	testutils.AssertEqual(t, "t2", tags[0].EntityID)
}

func getTagStore(tags []models.Tag) (tag_store.TagStore, int) {
	store := tag_store.NewMemoryTagStore()
	tagsAdded := 0

	for _, tag := range tags {
		tagsAdded++
		store.Insert(tag)
	}

	return store, tagsAdded
}
