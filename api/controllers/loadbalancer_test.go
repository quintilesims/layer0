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

func TestCreateLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb1",
		EnvironmentID:    "e1",
		IsPublic:         true,
		Ports:            []models.Port{},
		HealthCheck: models.HealthCheck{
			Target:             "80",
			Interval:           60,
			Timeout:            60,
			HealthyThreshold:   3,
			UnhealthyThreshold: 3,
		},
	}

	LoadBalancerModel := models.LoadBalancer{
		EnvironmentID:   "e1",
		EnvironmentName: "environment1",
		HealthCheck: models.HealthCheck{
			Target:             "80",
			Interval:           60,
			Timeout:            60,
			HealthyThreshold:   3,
			UnhealthyThreshold: 3,
		},
		IsPublic:         true,
		LoadBalancerID:   "lb1",
		LoadBalancerName: "loadbalancer1",
		Ports:            []models.Port{},
		ServiceID:        "s1",
		ServiceName:      "service1",
		URL:              "http://some-url.com",
	}

	mockJobScheduler := mock_scheduler.NewMockJobScheduler(ctrl)
	mockLoadBalancer := mock_provider.NewMockLoadBalancerProvider(ctrl)
	controller := NewLoadBalancerController(mockLoadBalancer, mockJobScheduler)

	mockLoadBalancer.EXPECT().
		Create(req).
		Return(&LoadBalancerModel, nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.CreateLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.LoadBalancer
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 202, recorder.Code)
	assert.Equal(t, LoadBalancerModel, response)
}

func TestDeleteLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobScheduler := scheduler.ScheduleJobFunc(func(req models.CreateJobRequest) (string, error) {
		assert.Equal(t, job.DeleteLoadBalancerJob, req.JobType)
		assert.Equal(t, "lb1", req.Request)

		return "j1", nil
	})

	mockLoadBalancer := mock_provider.NewMockLoadBalancerProvider(ctrl)
	controller := NewLoadBalancerController(mockLoadBalancer, mockJobScheduler)

	c := newFireballContext(t, nil, map[string]string{"id": "lb1"})
	resp, err := controller.DeleteLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 202, recorder.Code)
	assert.Equal(t, "j1", recorder.HeaderMap.Get("X-JobID"))
	assert.Equal(t, "/job/j1", recorder.HeaderMap.Get("Location"))
}

func TestGetLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	LoadBalancerModel := models.LoadBalancer{
		EnvironmentID:    "e1",
		EnvironmentName:  "environment1",
		HealthCheck:      models.HealthCheck{},
		IsPublic:         true,
		LoadBalancerID:   "lb1",
		LoadBalancerName: "loadbalancer1",
		Ports:            []models.Port{},
		ServiceID:        "s1",
		ServiceName:      "service1",
		URL:              "http://some-url.com",
	}

	mockJobScheduler := mock_scheduler.NewMockJobScheduler(ctrl)
	mockLoadBalancer := mock_provider.NewMockLoadBalancerProvider(ctrl)
	controller := NewLoadBalancerController(mockLoadBalancer, mockJobScheduler)

	mockLoadBalancer.EXPECT().
		Read("lb1").
		Return(&LoadBalancerModel, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "lb1"})
	resp, err := controller.GetLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.LoadBalancer
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, LoadBalancerModel, response)
}

func TestListLoadBalancers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	LoadBalancerSummaries := []models.LoadBalancerSummary{
		{
			LoadBalancerID:   "lb1",
			LoadBalancerName: "LoadBalancer1",
			EnvironmentID:    "e1",
			EnvironmentName:  "environment1",
		},
		{
			LoadBalancerID:   "lb2",
			LoadBalancerName: "LoadBalancer2",
			EnvironmentID:    "e2",
			EnvironmentName:  "environment2",
		},
	}

	mockJobScheduler := mock_scheduler.NewMockJobScheduler(ctrl)
	mockLoadBalancer := mock_provider.NewMockLoadBalancerProvider(ctrl)
	controller := NewLoadBalancerController(mockLoadBalancer, mockJobScheduler)

	mockLoadBalancer.EXPECT().
		List().
		Return(LoadBalancerSummaries, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.ListLoadBalancers(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.LoadBalancerSummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, LoadBalancerSummaries, response)
}
