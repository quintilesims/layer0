package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/api/scheduler"
	"github.com/quintilesims/layer0/api/scheduler/mock_scheduler"
	"github.com/quintilesims/layer0/common/job"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := models.CreateEnvironmentRequest{
		EnvironmentName: "env",
		InstanceSize:    "m3.medium",
		MinClusterCount: 1,
		OperatingSystem: "linux",
		AMIID:           "ami123",
	}

	environmentModel := models.Environment{
		EnvironmentID:   "e1",
		EnvironmentName: "env",
		InstanceSize:    "m3.medium",
		ClusterCount:    1,
		SecurityGroupID: "sg1",
		OperatingSystem: "linux",
		AMIID:           "ami123",
		Links:           []string{"e2"},
	}

	mockJobScheduler := mock_scheduler.NewMockJobScheduler(ctrl)
	mockEnvironment := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironment, mockJobScheduler)

	mockEnvironment.EXPECT().
		Create(req).
		Return(&environmentModel, nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.CreateEnvironment(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Environment
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 202, recorder.Code)
	assert.Equal(t, environmentModel, response)
}

func TestDeleteEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobScheduler := scheduler.ScheduleJobFunc(func(req models.CreateJobRequest) (string, error) {
		assert.Equal(t, job.DeleteEnvironmentJob, req.JobType)
		assert.Equal(t, "e1", req.Request)

		return "j1", nil
	})

	mockEnvironment := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironment, mockJobScheduler)

	c := newFireballContext(t, nil, map[string]string{"id": "e1"})
	resp, err := controller.DeleteEnvironment(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 202, recorder.Code)
	assert.Equal(t, "j1", recorder.HeaderMap.Get("X-JobID"))
	assert.Equal(t, "/job/j1", recorder.HeaderMap.Get("Location"))
}

func TestGetEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	environmentModel := models.Environment{
		EnvironmentID:   "e1",
		EnvironmentName: "env",
		InstanceSize:    "m3.medium",
		ClusterCount:    1,
		SecurityGroupID: "sg1",
		OperatingSystem: "linux",
		AMIID:           "ami123",
		Links:           []string{"e2"},
	}

	mockJobScheduler := mock_scheduler.NewMockJobScheduler(ctrl)
	mockEnvironment := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironment, mockJobScheduler)

	mockEnvironment.EXPECT().
		Read("e1").
		Return(&environmentModel, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "e1"})
	resp, err := controller.GetEnvironment(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Environment
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, environmentModel, response)
}

func TestListEnvironments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	environmentSummaries := []models.EnvironmentSummary{
		{
			EnvironmentID:   "e1",
			EnvironmentName: "env1",
			OperatingSystem: "linux",
		},
		{
			EnvironmentID:   "e2",
			EnvironmentName: "env2",
			OperatingSystem: "windows",
		},
	}

	mockJobScheduler := mock_scheduler.NewMockJobScheduler(ctrl)
	mockEnvironment := mock_provider.NewMockEnvironmentProvider(ctrl)
	controller := NewEnvironmentController(mockEnvironment, mockJobScheduler)

	mockEnvironment.EXPECT().
		List().
		Return(environmentSummaries, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.ListEnvironments(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.EnvironmentSummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, environmentSummaries, response)
}
