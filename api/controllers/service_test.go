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

func TestCreateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := models.CreateServiceRequest{
		DeployID:       "deploy_id",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		ServiceName:    "service_name",
	}

	serviceModel := models.Service{
		Deployments:      []models.Deployment{},
		DesiredCount:     2,
		EnvironmentID:    "env_id",
		EnvironmentName:  "env_name",
		LoadBalancerID:   "lb_id",
		LoadBalancerName: "lb_name",
		PendingCount:     2,
		RunningCount:     1,
		ServiceID:        "service_id",
		ServiceName:      "service_name",
	}

	mockService := mock_provider.NewMockServiceProvider(ctrl)
	mockJobScheduler := mock_scheduler.NewMockJobScheduler(ctrl)
	controller := NewServiceController(mockService, mockJobScheduler)

	mockService.EXPECT().
		Create(req).
		Return(&serviceModel, nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.CreateService(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Service
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 202, recorder.Code)
	assert.Equal(t, serviceModel, response)
}

func TestDeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobScheduler := scheduler.ScheduleJobFunc(func(req models.CreateJobRequest) (string, error) {
		assert.Equal(t, job.DeleteServiceJob, req.JobType)
		assert.Equal(t, "service_id", req.Request)

		return "j1", nil
	})

	mockService := mock_provider.NewMockServiceProvider(ctrl)
	controller := NewServiceController(mockService, mockJobScheduler)

	c := newFireballContext(t, nil, map[string]string{"id": "service_id"})
	resp, err := controller.DeleteService(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 202, recorder.Code)
	assert.Equal(t, "j1", recorder.HeaderMap.Get("X-JobID"))
	assert.Equal(t, "/job/j1", recorder.HeaderMap.Get("Location"))
}

func TestGetService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceModel := models.Service{
		Deployments:      ([]models.Deployment(nil)),
		DesiredCount:     2,
		EnvironmentID:    "env_id",
		EnvironmentName:  "env_name",
		LoadBalancerID:   "lb_id",
		LoadBalancerName: "lb_name",
		PendingCount:     2,
		RunningCount:     1,
		ServiceID:        "service_id",
		ServiceName:      "service_name",
	}

	mockService := mock_provider.NewMockServiceProvider(ctrl)
	mockJobScheduler := mock_scheduler.NewMockJobScheduler(ctrl)
	controller := NewServiceController(mockService, mockJobScheduler)

	mockService.EXPECT().
		Read("service_id").
		Return(&serviceModel, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "service_id"})
	resp, err := controller.GetService(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Service
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, serviceModel, response)
}

func TestListServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceSummaries := []models.ServiceSummary{
		{
			ServiceID:       "service_id",
			ServiceName:     "service_name",
			EnvironmentID:   "env_id",
			EnvironmentName: "env_name",
		},
		{
			ServiceID:       "service_id",
			ServiceName:     "service_name",
			EnvironmentID:   "env_id",
			EnvironmentName: "env_name",
		},
	}

	mockService := mock_provider.NewMockServiceProvider(ctrl)
	mockJobScheduler := mock_scheduler.NewMockJobScheduler(ctrl)
	controller := NewServiceController(mockService, mockJobScheduler)

	mockService.EXPECT().
		List().
		Return(serviceSummaries, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.ListServices(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.ServiceSummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, serviceSummaries, response)
}
