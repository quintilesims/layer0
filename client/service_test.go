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
		ServiceName:    "name",
		EnvironmentID:  "eid",
		DeployID:       "did",
		LoadBalancerID: "lid",
		Scale:          3,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/service")

		var body models.CreateServiceRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "jid"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.CreateService(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "jid", jobID)
}

func TestDeleteService(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/service/sid")

		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "jid"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.DeleteService("sid")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, jobID, "jid")
}

func TestListServices(t *testing.T) {
	expected := []*models.ServiceSummary{
		{ServiceID: "sid1"},
		{ServiceID: "sid2"},
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
	expected := &models.Service{
		ServiceID:   "sid",
		ServiceName: "ename",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/service/sid")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadService("sid")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestReadServiceLogs(t *testing.T) {
	expected := []*models.LogFile{
		{ContainerName: "c1"},
		{ContainerName: "c2"},
	}

	query := url.Values{}
	query.Set(LogQueryParamTail, "100")
	query.Set(LogQueryParamStart, "2000-01-01 00:00")
	query.Set(LogQueryParamEnd, "2000-01-01 12:12")

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/service/sid/logs")
		assert.Equal(t, query, r.URL.Query())

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadServiceLogs("sid", query)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestUpdateService(t *testing.T) {
	deployID := "did"
	scale := 1
	req := models.UpdateServiceRequest{
		DeployID: &deployID,
		Scale:    &scale,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PATCH")
		assert.Equal(t, r.URL.Path, "/service/sid")

		var body models.UpdateServiceRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "jid"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.UpdateService("sid", req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "jid", jobID)
}
