package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoadBalancerProvider := mock_provider.NewMockLoadBalancerProvider(ctrl)
	controller := NewLoadBalancerController(mockLoadBalancerProvider)

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		IsPublic:         true,
		Ports: []models.Port{
			{HostPort: 443, ContainerPort: 80, Protocol: "https", CertificateName: "cert"},
			{HostPort: 22, ContainerPort: 22, Protocol: "tcp"},
		},
		HealthCheck: models.HealthCheck{
			Target:             "tcp:80",
			Interval:           5,
			Timeout:            6,
			HealthyThreshold:   7,
			UnhealthyThreshold: 8,
		},
	}

	mockLoadBalancerProvider.EXPECT().
		Create(req).
		Return("lb_id", nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.createLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.CreateEntityResponse
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "lb_id", response.EntityID)
}

func TestDeleteLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoadBalancerProvider := mock_provider.NewMockLoadBalancerProvider(ctrl)
	controller := NewLoadBalancerController(mockLoadBalancerProvider)

	mockLoadBalancerProvider.EXPECT().
		Delete("lb_id").
		Return(nil)

	c := newFireballContext(t, nil, map[string]string{"id": "lb_id"})
	resp, err := controller.deleteLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)
}

func TestListLoadBalancers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoadBalancerProvider := mock_provider.NewMockLoadBalancerProvider(ctrl)
	controller := NewLoadBalancerController(mockLoadBalancerProvider)

	expected := []models.LoadBalancerSummary{
		{
			LoadBalancerID:   "lb_id1",
			LoadBalancerName: "lb_name1",
			EnvironmentID:    "env_id1",
			EnvironmentName:  "env_name1",
		},
		{
			LoadBalancerID:   "lb_id2",
			LoadBalancerName: "lbd_name2",
			EnvironmentID:    "env_id2",
			EnvironmentName:  "env_name2",
		},
	}

	mockLoadBalancerProvider.EXPECT().
		List().
		Return(expected, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.listLoadBalancers(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.LoadBalancerSummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestReadLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := models.LoadBalancer{
		LoadBalancerID:   "lb_id",
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		EnvironmentName:  "env_name",
		ServiceID:        "svc_id",
		ServiceName:      "svc_name",
		IsPublic:         true,
		URL:              "url",
		Ports: []models.Port{
			{HostPort: 443, ContainerPort: 80, Protocol: "https", CertificateName: "cert"},
			{HostPort: 22, ContainerPort: 22, Protocol: "tcp"},
		},
		HealthCheck: models.HealthCheck{
			Target:             "tcp:80",
			Interval:           1,
			Timeout:            2,
			HealthyThreshold:   3,
			UnhealthyThreshold: 4,
		},
	}

	mockLoadBalancerProvider := mock_provider.NewMockLoadBalancerProvider(ctrl)
	controller := NewLoadBalancerController(mockLoadBalancerProvider)

	mockLoadBalancerProvider.EXPECT().
		Read("lb_id").
		Return(&expected, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "lb_id"})
	resp, err := controller.readLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.LoadBalancer
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestUpdateLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoadBalancerProvider := mock_provider.NewMockLoadBalancerProvider(ctrl)
	controller := NewLoadBalancerController(mockLoadBalancerProvider)

	ports := []models.Port{
		{HostPort: 443, ContainerPort: 80, Protocol: "https", CertificateName: "cert"},
		{HostPort: 22, ContainerPort: 22, Protocol: "tcp"},
	}

	healthCheck := models.HealthCheck{
		Target:             "tcp:80",
		Interval:           1,
		Timeout:            2,
		HealthyThreshold:   3,
		UnhealthyThreshold: 4,
	}

	req := models.UpdateLoadBalancerRequest{
		Ports:       &ports,
		HealthCheck: &healthCheck,
	}

	mockLoadBalancerProvider.EXPECT().
		Update("lb_id", req).
		Return(nil)

	c := newFireballContext(t, req, map[string]string{"id": "lb_id"})
	resp, err := controller.updateLoadBalancer(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)
}
