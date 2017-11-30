package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job/mock_job"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoadBalancerProvider := mock_provider.NewMockLoadBalancerProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewLoadBalancerController(mockLoadBalancerProvider, mockJobStore, tagStore)

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

	mockJobStore.EXPECT().
		Insert(models.CreateLoadBalancerJob, gomock.Any()).
		Return("jid", nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.CreateLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "jid", response.JobID)
}

func TestDeleteLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoadBalancerProvider := mock_provider.NewMockLoadBalancerProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewLoadBalancerController(mockLoadBalancerProvider, mockJobStore, tagStore)

	mockJobStore.EXPECT().
		Insert(models.DeleteLoadBalancerJob, "lid").
		Return("jid", nil)

	c := newFireballContext(t, nil, map[string]string{"id": "lid"})
	resp, err := controller.DeleteLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "jid", response.JobID)
}

func TestReadLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoadBalancerProvider := mock_provider.NewMockLoadBalancerProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewLoadBalancerController(mockLoadBalancerProvider, mockJobStore, tagStore)

	loadBalancerModel := models.LoadBalancer{
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

	mockLoadBalancerProvider.EXPECT().
		Read("lb1").
		Return(&loadBalancerModel, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "lb1"})
	resp, err := controller.ReadLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.LoadBalancer
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, loadBalancerModel, response)
}

func TestListLoadBalancers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoadBalancerProvider := mock_provider.NewMockLoadBalancerProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewLoadBalancerController(mockLoadBalancerProvider, mockJobStore, tagStore)

	loadBalancerSummaries := []models.LoadBalancerSummary{
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

	mockLoadBalancerProvider.EXPECT().
		List().
		Return(loadBalancerSummaries, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.ListLoadBalancers(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.LoadBalancerSummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, loadBalancerSummaries, response)
}

func TestUpdateLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoadBalancerProvider := mock_provider.NewMockLoadBalancerProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewLoadBalancerController(mockLoadBalancerProvider, mockJobStore, tagStore)

	req := models.UpdateLoadBalancerRequest{
		Ports: &[]models.Port{},
		HealthCheck: &models.HealthCheck{
			Target:             "80",
			Interval:           60,
			Timeout:            60,
			HealthyThreshold:   3,
			UnhealthyThreshold: 3,
		},
	}

	mockJobStore.EXPECT().
		Insert(models.UpdateLoadBalancerJob, gomock.Any()).
		Return("jid", nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.UpdateLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Job
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "jid", response.JobID)
}
