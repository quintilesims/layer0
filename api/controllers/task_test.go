package controllers

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	controller := NewTaskController(mockTaskProvider)

	req := models.CreateTaskRequest{
		TaskName:      "tsk_name",
		EnvironmentID: "env_id",
		DeployID:      "dpl_id",
		ContainerOverrides: []models.ContainerOverride{
			{ContainerName: "c1", EnvironmentOverrides: map[string]string{"k1": "v1"}},
			{ContainerName: "c2", EnvironmentOverrides: map[string]string{"k2": "v2"}},
		},
	}

	mockTaskProvider.EXPECT().
		Create(req).
		Return("tsk_id", nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.createTask(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.CreateEntityResponse
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "tsk_id", response.EntityID)
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	controller := NewTaskController(mockTaskProvider)

	mockTaskProvider.EXPECT().
		Delete("tsk_id").
		Return(nil)

	c := newFireballContext(t, nil, map[string]string{"id": "tsk_id"})
	resp, err := controller.deleteTask(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)
}

func TestListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	controller := NewTaskController(mockTaskProvider)

	expected := []models.TaskSummary{
		{
			TaskID:          "tsk_id1",
			TaskName:        "tsk_name1",
			EnvironmentID:   "env_id1",
			EnvironmentName: "env_name1",
		},
		{
			TaskID:          "tsk_id2",
			TaskName:        "tskd_name2",
			EnvironmentID:   "env_id2",
			EnvironmentName: "env_name2",
		},
	}

	mockTaskProvider.EXPECT().
		List().
		Return(expected, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.listTasks(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.TaskSummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestReadTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := models.Task{
		TaskID:          "tsk_id",
		TaskName:        "tsk_name",
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		DeployID:        "dpl_id",
		DeployName:      "dpl_name",
		DeployVersion:   "1",
		Status:          "RUNNING",
		Containers: []models.Container{
			{ContainerName: "c1", Status: "RUNNING", ExitCode: 0},
			{ContainerName: "c2", Status: "STOPPED", ExitCode: 1},
		},
	}

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	controller := NewTaskController(mockTaskProvider)

	mockTaskProvider.EXPECT().
		Read("tsk_id").
		Return(&expected, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "tsk_id"})
	resp, err := controller.readTask(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Task
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestReadTaskLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	controller := NewTaskController(mockTaskProvider)

	expected := []models.LogFile{
		{
			ContainerName: "apline",
			Lines:         []string{"hello", "world"},
		},
	}

	tail := "100"
	start, err := time.Parse(client.TimeLayout, "2001-01-02 10:00")
	if err != nil {
		t.Fatalf("Failed to parse start: %v", err)
	}

	end, err := time.Parse(client.TimeLayout, "2001-01-02 12:00")
	if err != nil {
		t.Fatalf("Failed to parse end: %v", err)
	}

	mockTaskProvider.EXPECT().
		Logs("tsk_id", 100, start, end).
		Return(expected, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "tsk_id"})
	c.Request.URL.RawQuery = fmt.Sprintf("tail=%s&start=%s&end=%s",
		tail,
		start.Format(client.TimeLayout),
		end.Format(client.TimeLayout))

	resp, err := controller.readTaskLogs(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.LogFile
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}
