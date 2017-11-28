package controllers

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job/mock_job"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewServiceController(mockServiceProvider, mockJobStore, tagStore)

	req := models.CreateServiceRequest{
		DeployID:       "deploy_id",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		ServiceName:    "service_name",
	}

	mockJobStore.EXPECT().
		Insert(models.CreateServiceJob, gomock.Any()).
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
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewServiceController(mockServiceProvider, mockJobStore, tagStore)

	mockJobStore.EXPECT().
		Insert(models.DeleteServiceJob, "sid").
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
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewServiceController(mockServiceProvider, mockJobStore, tagStore)

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

func TestGetServiceLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewServiceController(mockServiceProvider, mockJobStore, tagStore)

	logFiles := []models.LogFile{
		{
			ContainerName: "apline",
			Lines:         []string{"hello", "world"},
		},
	}

	tail := "100"
	start, err := time.Parse(TIME_LAYOUT, "2001-01-02 10:00")
	if err != nil {
		t.Fatalf("Failed to parse start: %v", err)
	}

	end, err := time.Parse(TIME_LAYOUT, "2001-01-02 12:00")
	if err != nil {
		t.Fatalf("Failed to parse end: %v", err)
	}

	mockServiceProvider.EXPECT().
		Logs("service_id", 100, start, end).
		Return(logFiles, nil)

	c := newFireballContext(t, nil, map[string]string{"id": "service_id"})
	c.Request.URL.RawQuery = fmt.Sprintf("tail=%s&start=%s&end=%s",
		tail,
		start.Format(TIME_LAYOUT),
		end.Format(TIME_LAYOUT))

	resp, err := controller.GetServiceLogs(c)
	if err != nil {
		t.Fatal(err)
	}

	var response []models.LogFile
	recorder := unmarshalBody(t, resp, &response)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, logFiles, response)
}

func TestListServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	tagStore := tag.NewMemoryStore()
	controller := NewServiceController(mockServiceProvider, mockJobStore, tagStore)

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
