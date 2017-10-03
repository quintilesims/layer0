package controllers

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/job/mock_job"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	controller := NewTaskController(mockTaskProvider, mockJobStore)

	req := models.CreateTaskRequest{
		DeployID:      "deploy_id",
		EnvironmentID: "env_id",
		TaskName:      "task_name",
	}

	mockJobStore.EXPECT().
		Insert(job.CreateTaskJob, gomock.Any()).
		Return("jid", nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.CreateTask(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "jid", response.JobID)
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	controller := NewTaskController(mockTaskProvider, mockJobStore)

	mockJobStore.EXPECT().
		Insert(job.DeleteTaskJob, "tid").
		Return("jid", nil)

	c := newFireballContext(t, nil, map[string]string{"id": "tid"})
	resp, err := controller.DeleteTask(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "jid", response.JobID)
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	controller := NewTaskController(mockTaskProvider, mockJobStore)

	taskModel := models.Task{
		TaskID:          "task_id",
		TaskName:        "task_name",
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		DeployID:        "deploy_id",
		DeployName:      "deploy_name",
		DeployVersion:   "5",
		Status:          "RUNNING",
		Containers: []models.Container{
			{
				ContainerName: "name",
				Status:        "RUNNING",
				ExitCode:      0,
				Meta:          "",
			},
		},
	}

	mockTaskProvider.EXPECT().
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

func TestGetTaskLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	controller := NewTaskController(mockTaskProvider, mockJobStore)

	logFiles := []models.LogFile{
		{
			ContainerName: "alpine",
			Lines:         []string{"hello", "world"},
		},
	}

	tail := "100"
	start, err := time.Parse(TIME_LAYOUT, "2001-01-02 10:00")
	if err != nil {
		t.Fatalf("Failed to parse start: %v", err)
	}

	end, err := time.Parse(TIME_LAYOUT, "2001-01-02 12:00")
	if err != nil {
		t.Fatalf("Failed to parse end: %v", err)
	}

	mockTaskProvider.EXPECT().
		Logs("task_id", 100, start, end).
		Return(logFiles, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "task_id"})
	c.Request.URL.RawQuery = fmt.Sprintf("tail=%s&start=%s&end=%s",
		tail,
		start.Format(TIME_LAYOUT),
		end.Format(TIME_LAYOUT))

	resp, err := controller.GetTaskLogs(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.LogFile
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, logFiles, response)
}

func TestListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	controller := NewTaskController(mockTaskProvider, mockJobStore)

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

	mockTaskProvider.EXPECT().
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
