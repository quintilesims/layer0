package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeploy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := models.CreateDeployRequest{
		DeployName: "deploy1",
		Dockerrun:  ([]byte("content")),
	}

	DeployModel := models.Deploy{
		DeployID:   "d1",
		DeployName: "deploy1",
		Dockerrun:  ([]byte("content")),
		Version:    "1",
	}

	mockDeploy := mock_provider.NewMockDeployProvider(ctrl)
	controller := NewDeployController(mockDeploy)

	mockDeploy.EXPECT().
		Create(req).
		Return(&DeployModel, nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.CreateDeploy(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Deploy
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 202, recorder.Code)
	assert.Equal(t, DeployModel, response)
}

func TestDeleteDeploy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeploy := mock_provider.NewMockDeployProvider(ctrl)
	controller := NewDeployController(mockDeploy)

	mockDeploy.EXPECT().
		Delete("d1").
		Return(nil)

	c := newFireballContext(t, nil, map[string]string{"id": "d1"})
	resp, err := controller.DeleteDeploy(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Deploy
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
}

func TestGetDeploy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	DeployModel := models.Deploy{
		DeployID:   "d1",
		DeployName: "deploy1",
		Dockerrun:  ([]byte("content")),
		Version:    "1",
	}

	mockDeploy := mock_provider.NewMockDeployProvider(ctrl)
	controller := NewDeployController(mockDeploy)

	mockDeploy.EXPECT().
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

	DeploySummaries := []models.DeploySummary{
		{
			DeployID:   "d1",
			DeployName: "deploy1",
			Version:    "1",
		},
		{
			DeployID:   "d2",
			DeployName: "deploy2",
			Version:    "1",
		},
	}

	mockDeploy := mock_provider.NewMockDeployProvider(ctrl)
	controller := NewDeployController(mockDeploy)

	mockDeploy.EXPECT().
		List().
		Return(DeploySummaries, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.ListDeploys(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.DeploySummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, DeploySummaries, response)
}
