package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	mockJobScheduler := mock_job.NewMockScheduler(ctrl)
	controller := NewServiceController(mockServiceProvider, mockJobScheduler)

	req := models.CreateServiceRequest{
		DeployID:       "deploy_id",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		ServiceName:    "service_name",
	}

	sjr := models.ScheduleJobRequest{
		JobType: job.CreateServiceJob.String(),
		Request: req,
	}

	mockJobScheduler.EXPECT().
		Schedule(sjr).
		Return("jid", nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.CreateService(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "jid", response.JobID)
}

func TestDeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	mockJobScheduler := mock_job.NewMockScheduler(ctrl)
	controller := NewServiceController(mockServiceProvider, mockJobScheduler)

	sjr := models.ScheduleJobRequest{
		JobType: job.DeleteServiceJob.String(),
		Request: "sid",
	}

	mockJobScheduler.EXPECT().
		Schedule(sjr).
		Return("jid", nil)

	c := newFireballContext(t, nil, map[string]string{"id": "sid"})
	resp, err := controller.DeleteService(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "jid", response.JobID)
}

func TestGetService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	mockJobScheduler := mock_job.NewMockScheduler(ctrl)
	controller := NewServiceController(mockServiceProvider, mockJobScheduler)

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

	mockServiceProvider.EXPECT().
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

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	mockJobScheduler := mock_job.NewMockScheduler(ctrl)
	controller := NewServiceController(mockServiceProvider, mockJobScheduler)

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

	mockServiceProvider.EXPECT().
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
