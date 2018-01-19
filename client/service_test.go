package client

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateService(t *testing.T) {
	req := models.CreateServiceRequest{
		ServiceName:    "svc_name",
		EnvironmentID:  "env_id",
		DeployID:       "dpl_id",
		LoadBalancerID: "lb_id",
		Scale:          3,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/service")

		var body models.CreateServiceRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateEntityResponse{EntityID: "svc_id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	serviceID, err := client.CreateService(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "svc_id", serviceID)
}

func TestDeleteService(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/service/svc_id")

		MarshalAndWrite(t, w, nil, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteService("svc_id"); err != nil {
		t.Fatal(err)
	}
}

func TestListServices(t *testing.T) {
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

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/service")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ListServices()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestReadService(t *testing.T) {
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

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/service/svc_id")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadService("svc_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestReadServiceLogs(t *testing.T) {
	expected := []models.LogFile{
		{
			ContainerName: "apline",
			Lines:         []string{"hello", "world"},
		},
	}

	query := url.Values{}
	query.Set(LogQueryParamTail, "100")
	query.Set(LogQueryParamStart, "2000-01-01 00:00")
	query.Set(LogQueryParamEnd, "2000-01-01 12:12")

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/service/svc_id/logs")
		assert.Equal(t, query, r.URL.Query())

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadServiceLogs("svc_id", query)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestUpdateService(t *testing.T) {
	deployID := "dpl_id"
	scale := 1

	req := models.UpdateServiceRequest{
		DeployID: &deployID,
		Scale:    &scale,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PATCH")
		assert.Equal(t, r.URL.Path, "/service/svc_id")

		var body models.UpdateServiceRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, nil, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.UpdateService("svc_id", req); err != nil {
		t.Fatal(err)
	}
}
