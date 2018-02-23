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

func TestCreateEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironmentProvider)

	req := models.CreateEnvironmentRequest{
		EnvironmentName: "env",
		InstanceType:    "t2.small",
		Scale:           3,
		OperatingSystem: "linux",
		AMIID:           "ami123",
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
		CurrentScale:    2,
		DesiredScale:    3,
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

func TestReadEnvironmentLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironmentProvider)

	expected := []models.LogFile{
		{
			ContainerName: "alpine",
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

	mockEnvironmentProvider.EXPECT().
		Logs("env_id", 100, start, end).
		Return(expected, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "env_id"})
	c.Request.URL.RawQuery = fmt.Sprintf("tail=%s&start=%s&end=%s",
		tail,
		start.Format(client.TimeLayout),
		end.Format(client.TimeLayout))

	resp, err := controller.readEnvironmentLogs(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.LogFile
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestUpdateEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironmentProvider)

	scale := 1
	links := []string{"link1", "link2"}
	req := models.UpdateEnvironmentRequest{
		Scale: &scale,
		Links: &links,
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
