package logic

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTaskPopulateModel(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "t1", EntityType: "task", Key: "name", Value: "tsk_1"},
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env_1"},
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl_1"},

		{EntityID: "t2", EntityType: "task", Key: "name", Value: "tsk_2"},
		{EntityID: "t2", EntityType: "task", Key: "environment_id", Value: "e2"},
		{EntityID: "t2", EntityType: "task", Key: "deploy_id", Value: "d2"},
	})

	cases := map[*models.Task]func(*models.Task){
		{
			TaskID:        "t1",
			EnvironmentID: "e1",
			DeployID:      "d1",
		}: func(m *models.Task) {
			testutils.AssertEqual(t, m.TaskName, "tsk_1")
			testutils.AssertEqual(t, m.EnvironmentName, "env_1")
			testutils.AssertEqual(t, m.DeployName, "dpl_1")
		},
		{
			TaskID: "t2",
		}: func(m *models.Task) {
			testutils.AssertEqual(t, m.TaskName, "tsk_2")
			testutils.AssertEqual(t, m.EnvironmentID, "e2")
			testutils.AssertEqual(t, m.DeployID, "d2")
		},
	}

	taskLogic := NewL0TaskLogic(testLogic.Logic())
	for model, fn := range cases {
		if err := taskLogic.populateModel(model); err != nil {
			t.Fatal(err)
		}

		fn(model)
	}
}

func TestGetTask(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		GetTask("e1", "t1").
		Return(&models.Task{TaskID: "t1"}, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "t1", EntityType: "task", Key: "environment_id", Value: "e1"},
	})

	taskLogic := NewL0TaskLogic(testLogic.Logic())
	task, err := taskLogic.GetTask("t1")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, task.TaskID, "t1")
	testutils.AssertEqual(t, task.EnvironmentID, "e1")
}

func TestListTasks(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	taskARNs := []string{
		"arn1",
		"arn2",
		"extra",
	}

	testLogic.Backend.EXPECT().
		ListTasks().
		Return(taskARNs, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "env_id1", EntityType: "environment", Key: "name", Value: "env_name1"},
		{EntityID: "env_id2", EntityType: "environment", Key: "name", Value: "env_name2"},
		{EntityID: "tsk_id1", EntityType: "task", Key: "name", Value: "tsk_name1"},
		{EntityID: "tsk_id1", EntityType: "task", Key: "environment_id", Value: "env_id1"},
		{EntityID: "tsk_id1", EntityType: "task", Key: "arn", Value: "arn1"},
		{EntityID: "tsk_id2", EntityType: "task", Key: "name", Value: "tsk_name2"},
		{EntityID: "tsk_id2", EntityType: "task", Key: "environment_id", Value: "env_id2"},
		{EntityID: "tsk_id2", EntityType: "task", Key: "arn", Value: "arn2"},
	})

	taskLogic := NewL0TaskLogic(testLogic.Logic())
	result, err := taskLogic.ListTasks()
	if err != nil {
		t.Fatal(err)
	}

	expected := []*models.TaskSummary{
		{EnvironmentID: "env_id1", EnvironmentName: "env_name1", TaskID: "tsk_id1", TaskName: "tsk_name1"},
		{EnvironmentID: "env_id2", EnvironmentName: "env_name2", TaskID: "tsk_id2", TaskName: "tsk_name2"},
	}

	assert.Equal(t, expected, result)
}

func TestDeleteTask(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		DeleteTask("e1", "t1").
		Return(nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "t1", EntityType: "task", Key: "name", Value: "tsk"},
		{EntityID: "t1", EntityType: "task", Key: "environment_id", Value: "e1"},
		{EntityID: "extra", EntityType: "task", Key: "name", Value: "extra"},
	})

	taskLogic := NewL0TaskLogic(testLogic.Logic())
	if err := taskLogic.DeleteTask("t1"); err != nil {
		t.Fatal(err)
	}

	tags, err := testLogic.TagStore.SelectByType("task")
	if err != nil {
		t.Fatal(err)
	}

	// make sure the 'extra' tag is the only one left
	testutils.AssertEqual(t, len(tags), 1)
}

func TestCreateTask(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		CreateTask("e1", "name", "d1", nil).
		Return(&models.Task{TaskID: "t1"}, nil)

	request := models.CreateTaskRequest{
		TaskName:      "name",
		EnvironmentID: "e1",
		DeployID:      "d1",
		Copies:        2,
	}

	taskLogic := NewL0TaskLogic(testLogic.Logic())
	task, err := taskLogic.CreateTask(request)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, task.TaskID, "t1")
	testutils.AssertEqual(t, task.EnvironmentID, "e1")

	testLogic.AssertTagExists(t, models.Tag{EntityID: "t1", EntityType: "task", Key: "name", Value: "name"})
	testLogic.AssertTagExists(t, models.Tag{EntityID: "t1", EntityType: "task", Key: "environment_id", Value: "e1"})
	testLogic.AssertTagExists(t, models.Tag{EntityID: "t1", EntityType: "task", Key: "deploy_id", Value: "d1"})
}

func TestCreateTaskError_missingRequiredParams(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	taskLogic := NewL0TaskLogic(testLogic.Logic())

	cases := map[string]models.CreateTaskRequest{
		"Missing EnvironmentID": {
			TaskName: "name",
			DeployID: "d1",
		},
		"Missing TaskName": {
			EnvironmentID: "e1",
			DeployID:      "d1",
		},
		"Missing DeployID": {
			EnvironmentID: "e1",
			TaskName:      "name",
		},
	}

	for name, request := range cases {
		if _, err := taskLogic.CreateTask(request); err == nil {
			t.Errorf("Case %s: error was nil!", name)
		}
	}
}

func TestGetTaskLogs(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	logs := []*models.LogFile{
		{Name: "alpha", Lines: []string{"first", "second"}},
		{Name: "beta", Lines: []string{"first", "second", "third"}},
	}

	testLogic.Backend.EXPECT().
		GetTaskLogs("e1", "t1", "start", "end", 100).
		Return(logs, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "t1", EntityType: "task", Key: "environment_id", Value: "e1"},
	})

	taskLogic := NewL0TaskLogic(testLogic.Logic())
	received, err := taskLogic.GetTaskLogs("t1", "start", "end", 100)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, received, logs)
}
