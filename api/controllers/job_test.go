package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/entity/mock_entity"
	"github.com/quintilesims/layer0/common/job"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDeleteJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJob := mock_entity.NewMockJob(ctrl)
	mockJob.EXPECT().
		Delete().
		Return(nil)

	mockProvider := mock_entity.NewMockProvider(ctrl)
	mockProvider.EXPECT().
		GetJob("j1").
		Return(mockJob)

	controller := NewJobController(mockProvider)

	c := newFireballContext(t, nil, map[string]string{"id": "j1"})
	resp, err := controller.DeleteJob(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)
}

func TestGetJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobModel := models.Job{
		JobID:     "j1",
		JobStatus: job.InProgress,
		JobType:   job.DeleteEnvironmentJob,
		Request:   "e1",
	}

	mockJob := mock_entity.NewMockJob(ctrl)
	mockJob.EXPECT().
		Model().
		Return(&jobModel, nil)

	mockProvider := mock_entity.NewMockProvider(ctrl)
	mockProvider.EXPECT().
		GetJob("j1").
		Return(mockJob)

	controller := NewJobController(mockProvider)

	c := newFireballContext(t, nil, map[string]string{"id": "j1"})
	resp, err := controller.GetJob(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, jobModel, response)
}

func TestListJobs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobModels := []models.Job{
		{
			JobID:     "j1",
			JobStatus: job.InProgress,
			JobType:   job.DeleteEnvironmentJob,
			Request:   "e1",
		},
		{
			JobID:     "j2",
			JobStatus: job.Completed,
			JobType:   job.DeleteServiceJob,
			Request:   "s1",
		},
	}

	mockProvider := mock_entity.NewMockProvider(ctrl)
	mockProvider.EXPECT().
		ListJobIDs().
		Return([]string{"j1", "j2"}, nil)

	for i := range jobModels {
		model := jobModels[i]

		mockJob := mock_entity.NewMockJob(ctrl)
		mockJob.EXPECT().
			Model().
			Return(&model, nil)

		mockProvider.EXPECT().
			GetJob(model.JobID).
			Return(mockJob)
	}

	controller := NewJobController(mockProvider)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.ListJobs(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, jobModels, response)
}
