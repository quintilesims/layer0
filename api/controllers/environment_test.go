package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironmentProvider)

	req := models.CreateEnvironmentRequest{
		EnvironmentName:  "env_name",
		InstanceType:     "instance_type",
		UserDataTemplate: []byte("user_data_template"),
		MinScale:         1,
		MaxScale:         2,
		OperatingSystem:  "linux",
		AMIID:            "ami_id",
	}

	mockEnvironmentProvider.EXPECT().
		Create(req).
		Return("env_id", nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.createEnvironment(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.CreateEntityResponse
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "env_id", response.EntityID)
}

func TestDeleteEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironmentProvider)

	mockEnvironmentProvider.EXPECT().
		Delete("env_id").
		Return(nil)

	c := newFireballContext(t, nil, map[string]string{"id": "env_id"})
	resp, err := controller.deleteEnvironment(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)
}

func TestListEnvironments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironmentProvider)

	expected := []models.EnvironmentSummary{
		{
			EnvironmentID:   "env_id1",
			EnvironmentName: "env_name1",
			OperatingSystem: "linux",
		},
		{
			EnvironmentID:   "env_id2",
			EnvironmentName: "envd_name2",
			OperatingSystem: "windows",
		},
	}

	mockEnvironmentProvider.EXPECT().
		List().
		Return(expected, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.listEnvironments(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.EnvironmentSummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestReadEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := models.Environment{
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		MinScale:        1,
		CurrentScale:    2,
		MaxScale:        3,
		InstanceType:    "instance_type",
		SecurityGroupID: "security_group_id",
		OperatingSystem: "linux",
		AMIID:           "ami_id",
		Links:           []string{"link1", "link2"},
	}

	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironmentProvider)

	mockEnvironmentProvider.EXPECT().
		Read("env_id").
		Return(&expected, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "env_id"})
	resp, err := controller.readEnvironment(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Environment
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestUpdateEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironmentProvider)

	min := 1
	max := 2
	links := []string{"link1", "link2"}
	req := models.UpdateEnvironmentRequest{
		MinScale: &min,
		MaxScale: &max,
		Links:    &links,
	}

	mockEnvironmentProvider.EXPECT().
		Update("env_id", req).
		Return(nil)

	c := newFireballContext(t, req, map[string]string{"id": "env_id"})
	resp, err := controller.updateEnvironment(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)
}
