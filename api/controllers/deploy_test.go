package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/job/mock_job"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeploy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeployProvider := mock_provider.NewMockDeployProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewDeployController(mockDeployProvider, mockJobStore, tagStore)

	req := models.CreateDeployRequest{
		DeployName: "deploy1",
		DeployFile: []byte("content"),
	}

	mockJobStore.EXPECT().
		Insert(job.CreateDeployJob, gomock.Any()).
		Return("jid", nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.CreateDeploy(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "jid", response.JobID)
}

func TestDeleteDeploy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeployProvider := mock_provider.NewMockDeployProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewDeployController(mockDeployProvider, mockJobStore, tagStore)

	mockJobStore.EXPECT().
		Insert(job.DeleteDeployJob, "did").
		Return("jid", nil)

	c := newFireballContext(t, nil, map[string]string{"id": "did"})
	resp, err := controller.DeleteDeploy(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "jid", response.JobID)
}

func TestGetDeploy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	DeployModel := models.Deploy{
		DeployFile: []byte("content"),
		DeployID:   "d1",
		DeployName: "deploy1",
		Version:    "1",
	}

	mockDeployProvider := mock_provider.NewMockDeployProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewDeployController(mockDeployProvider, mockJobStore, tagStore)

	mockDeployProvider.EXPECT().
		Read("d1").
		Return(&DeployModel, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "d1"})
	resp, err := controller.GetDeploy(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Deploy
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, DeployModel, response)
}

func TestListDeploys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeployProvider := mock_provider.NewMockDeployProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewDeployController(mockDeployProvider, mockJobStore, tagStore)

	deploySummaries := []models.DeploySummary{
		{
			DeployID:   "d1",
			DeployName: "deploy1",
			Version:    "1",
		},
		{
			DeployID:   "d1",
			DeployName: "deploy1",
			Version:    "2",
		},
		{
			DeployID:   "d2",
			DeployName: "deploy2",
			Version:    "1",
		},
	}

	mockDeployProvider.EXPECT().
		List().
		Return(deploySummaries, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.ListDeploys(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.DeploySummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, deploySummaries, response)
}
