package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/job/mock_job"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDeleteJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobStore := mock_job.NewMockStore(ctrl)
	controller := NewJobController(mockJobStore)

	mockJobStore.EXPECT().
		Delete("jid").
		Return(nil)

	c := newFireballContext(t, nil, map[string]string{"id": "jid"})
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

	mockJobStore := mock_job.NewMockStore(ctrl)
	controller := NewJobController(mockJobStore)

	jobModel := models.Job{
		JobID:   "jid",
		Type:    job.CreateEnvironmentJob.String(),
		Status:  job.InProgress.String(),
		Request: "some data",
	}

	mockJobStore.EXPECT().
		SelectByID("jid").
		Return(&jobModel, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "jid"})
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

	mockJobStore := mock_job.NewMockStore(ctrl)
	controller := NewJobController(mockJobStore)

	jobModels := []*models.Job{
		{
			JobID:   "j1",
			Type:    job.CreateEnvironmentJob.String(),
			Status:  job.InProgress.String(),
			Request: "some data",
		},
		{
			JobID:   "j2",
			Type:    job.DeleteServiceJob.String(),
			Status:  job.Completed.String(),
			Request: "sid",
		},
	}

	mockJobStore.EXPECT().
		SelectAll().
		Return(jobModels, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.ListJobs(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []*models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, jobModels, response)
}
