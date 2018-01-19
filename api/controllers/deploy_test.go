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

	mockDeployProvider := mock_provider.NewMockDeployProvider(ctrl)
	controller := NewDeployController(mockDeployProvider)

	req := models.CreateDeployRequest{
		DeployName: "dpl_name",
		DeployFile: []byte("content"),
	}

	mockDeployProvider.EXPECT().
		Create(req).
		Return("dpl_id", nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.createDeploy(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.CreateEntityResponse
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "dpl_id", response.EntityID)
}

func TestDeleteDeploy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeployProvider := mock_provider.NewMockDeployProvider(ctrl)
	controller := NewDeployController(mockDeployProvider)

	mockDeployProvider.EXPECT().
		Delete("dpl_id").
		Return(nil)

	c := newFireballContext(t, nil, map[string]string{"id": "dpl_id"})
	resp, err := controller.deleteDeploy(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)
}

func TestListDeploys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeployProvider := mock_provider.NewMockDeployProvider(ctrl)
	controller := NewDeployController(mockDeployProvider)

	expected := []models.DeploySummary{
		{
			DeployID:   "dpl_id1",
			DeployName: "dpl_name1",
			Version:    "1",
		},
		{
			DeployID:   "dpl_id2",
			DeployName: "dpld_name2",
			Version:    "2",
		},
	}

	mockDeployProvider.EXPECT().
		List().
		Return(expected, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.listDeploys(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.DeploySummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestReadDeploy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := models.Deploy{
		DeployID:   "dpl_id",
		DeployName: "dpl_name",
		DeployFile: []byte("content"),
		Version:    "1",
	}

	mockDeployProvider := mock_provider.NewMockDeployProvider(ctrl)
	controller := NewDeployController(mockDeployProvider)

	mockDeployProvider.EXPECT().
		Read("dpl_id").
		Return(&expected, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "dpl_id"})
	resp, err := controller.readDeploy(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Deploy
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}
