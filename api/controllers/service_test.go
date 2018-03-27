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

func TestCreateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	controller := NewServiceController(mockServiceProvider)

	req := models.CreateServiceRequest{
		ServiceName:    "svc_name",
		EnvironmentID:  "env_id",
		DeployID:       "dpl_id",
		LoadBalancerID: "lb_id",
		Scale:          3,
	}

	mockServiceProvider.EXPECT().
		Create(req).
		Return("svc_id", nil)

	c := newFireballContext(t, req, nil)
	resp, err := controller.createService(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.CreateEntityResponse
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "svc_id", response.EntityID)
}

func TestDeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	controller := NewServiceController(mockServiceProvider)

	mockServiceProvider.EXPECT().
		Delete("svc_id").
		Return(nil)

	c := newFireballContext(t, nil, map[string]string{"id": "svc_id"})
	resp, err := controller.deleteService(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)
}

func TestListServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	controller := NewServiceController(mockServiceProvider)

	expected := []models.ServiceSummary{
		{
			ServiceID:       "svc_id1",
			ServiceName:     "svc_name1",
			EnvironmentID:   "env_id1",
			EnvironmentName: "env_name1",
		},
		{
			ServiceID:       "svc_id2",
			ServiceName:     "svcd_name2",
			EnvironmentID:   "env_id2",
			EnvironmentName: "env_name2",
		},
	}

	mockServiceProvider.EXPECT().
		List().
		Return(expected, nil)

	c := newFireballContext(t, nil, nil)
	resp, err := controller.listServices(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.ServiceSummary
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestReadService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := models.Service{
		ServiceID:        "svc_id",
		ServiceName:      "svc_name",
		EnvironmentID:    "env_id",
		EnvironmentName:  "env_name",
		LoadBalancerID:   "lb_id",
		LoadBalancerName: "lb_name",
		DesiredCount:     3,
		PendingCount:     2,
		RunningCount:     1,
		Deployments: []models.Deployment{
			{
				DeployID:      "dpl_id1",
				DeployName:    "dpl_name1",
				DeployVersion: "1",
				Status:        "RUNNING",
			},
			{
				DeployID:      "dpl_id2",
				DeployName:    "dpl_name2",
				DeployVersion: "2",
				Status:        "STOPPED",
			},
		},
	}

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	controller := NewServiceController(mockServiceProvider)

	mockServiceProvider.EXPECT().
		Read("svc_id").
		Return(&expected, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "svc_id"})
	resp, err := controller.readService(c)
	if err != nil {
		t.Fatal(err)
	}

	var response models.Service
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestReadServiceLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	controller := NewServiceController(mockServiceProvider)

	expected := []models.LogFile{
		{
			ContainerName: "apline",
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

	mockServiceProvider.EXPECT().
		Logs("svc_id", 100, start, end).
		Return(expected, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "svc_id"})
	c.Request.URL.RawQuery = fmt.Sprintf("tail=%s&start=%s&end=%s",
		tail,
		start.Format(client.TimeLayout),
		end.Format(client.TimeLayout))

	resp, err := controller.readServiceLogs(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.LogFile
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, expected, response)
}

func TestUpdateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	controller := NewServiceController(mockServiceProvider)

	deployID := "dpl_id"
	scale := 1

	req := models.UpdateServiceRequest{
		DeployID: &deployID,
		Scale:    &scale,
	}

	mockServiceProvider.EXPECT().
		Update("svc_id", req).
		Return(nil)

	c := newFireballContext(t, req, map[string]string{"id": "svc_id"})
	resp, err := controller.updateService(c)
	if err != nil {
		t.Fatal(err)
	}

	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)
}
