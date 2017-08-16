package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := models.CreateTaskRequest{
		ContainerOverrides: ([]models.ContainerOverride(nil)),
		Copies:             1,
		DeployID:           "deploy_id",
		EnvironmentID:      "env_id",
		TaskName:           "task_name",
	}

	taskModel := models.Task{
		Copies:          ([]models.TaskCopy(nil)),
		DeployID:        "deploy_id",
		DeployName:      "deploy_name",
		DeployVersion:   "deploy_version",
		DesiredCount:    1,
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		PendingCount:    1,
		RunningCount:    0,
		TaskID:          "task_id",
		TaskName:        "task_name",
	}

	mockTask := mock_provider.NewMockTaskProvider(ctrl)
	controller := NewTaskController(mockTask)

	mockTask.EXPECT().
		Create(req).
		Return(&taskModel, nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.CreateTask(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Task
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 202, recorder.Code)
	assert.Equal(t, taskModel, response)
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTask := mock_provider.NewMockTaskProvider(ctrl)
	controller := NewTaskController(mockTask)

	mockTask.EXPECT().
		Delete("d1").
		Return(nil)

	c := newFireballContext(t, nil, map[string]string{"id": "d1"})
	resp, err := controller.DeleteTask(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Task
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskModel := models.Task{
		Copies:          ([]models.TaskCopy(nil)),
		DeployID:        "deploy_id",
		DeployName:      "deploy_name",
		DeployVersion:   "deploy_version",
		DesiredCount:    1,
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		PendingCount:    1,
		RunningCount:    0,
		TaskID:          "task_id",
		TaskName:        "task_name",
	}

	mockTask := mock_provider.NewMockTaskProvider(ctrl)
	controller := NewTaskController(mockTask)

	mockTask.EXPECT().
		Read("task_id").
		Return(&taskModel, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "task_id"})
	resp, err := controller.GetTask(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Task
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, taskModel, response)
}

func TestListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskSummaries := []models.TaskSummary{
		{
			TaskID:          "task_id",
			TaskName:        "task_name",
			EnvironmentID:   "env_id",
			EnvironmentName: "env_name",
		},
		{
			TaskID:          "task_id",
			TaskName:        "task_name",
			EnvironmentID:   "env_id",
			EnvironmentName: "env_name",
		},
	}

	mockTask := mock_provider.NewMockTaskProvider(ctrl)
	controller := NewTaskController(mockTask)

	mockTask.EXPECT().
		List().
		Return(taskSummaries, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.ListTasks(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.TaskSummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, taskSummaries, response)
}
