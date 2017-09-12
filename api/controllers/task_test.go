package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	mockJobScheduler := mock_job.NewMockScheduler(ctrl)
	controller := NewTaskController(mockTaskProvider, mockJobScheduler)

	req := models.CreateTaskRequest{
		ContainerOverrides: []models.ContainerOverride{},
		Copies:             1,
		DeployID:           "deploy_id",
		EnvironmentID:      "env_id",
		TaskName:           "task_name",
	}

	sjr := models.ScheduleJobRequest{
		JobType: job.CreateTaskJob.String(),
		Request: req,
	}

	mockJobScheduler.EXPECT().
		Schedule(sjr).
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
	mockJobScheduler := mock_job.NewMockScheduler(ctrl)
	controller := NewTaskController(mockTaskProvider, mockJobScheduler)

	sjr := models.ScheduleJobRequest{
		JobType: job.DeleteTaskJob.String(),
		Request: "tid",
	}

	mockJobScheduler.EXPECT().
		Schedule(sjr).
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
	mockJobScheduler := mock_job.NewMockScheduler(ctrl)
	controller := NewTaskController(mockTaskProvider, mockJobScheduler)

	taskModel := models.Task{
		Copies:          []models.TaskCopy{},
		DeployID:        "deploy_id",
		DeployName:      "deploy_name",
		DeployVersion:   "deploy_version",
		DesiredCount:    2,
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		PendingCount:    2,
		RunningCount:    1,
		TaskID:          "task_id",
		TaskName:        "task_name",
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

func TestListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	mockJobScheduler := mock_job.NewMockScheduler(ctrl)
	controller := NewTaskController(mockTaskProvider, mockJobScheduler)

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
