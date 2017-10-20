package client

import (
	"net/http"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateEnvironment(t *testing.T) {
	req := models.CreateEnvironmentRequest{
		EnvironmentName:  "name",
		InstanceSize:     "size",
		MinClusterCount:  1,
		UserDataTemplate: []byte("user_data"),
		OperatingSystem:  "os",
		AMIID:            "ami",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/environment")

		var body models.CreateEnvironmentRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "jid"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.CreateEnvironment(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "jid", jobID)
}

func TestDeleteEnvironment(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/environment/eid")

		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "jid"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.DeleteEnvironment("eid")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, jobID, "jid")
}

func TestListEnvironments(t *testing.T) {
	expected := []*models.EnvironmentSummary{
		{EnvironmentID: "eid1"},
		{EnvironmentID: "eid2"},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/environment")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ListEnvironments()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestReadEnvironment(t *testing.T) {
	expected := &models.Environment{
		EnvironmentID:   "eid",
		EnvironmentName: "ename",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/environment/eid")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadEnvironment("eid")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestUpdateEnvironment(t *testing.T) {
	count := 1
	req := models.UpdateEnvironmentRequest{
		EnvironmentID:   "eid",
		MinClusterCount: &count,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/environment")

		var body models.UpdateEnvironmentRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "jid"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.UpdateEnvironment(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "jid", jobID)
}

func TestCreateEnvironmentLink(t *testing.T) {
	req := models.CreateEnvironmentLinkRequest{
		SourceEnvironmentID: "env_id1",
		DestEnvironmentID:   "env_id2",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/environment/env_id1/link")

		var body models.CreateEnvironmentLinkRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "job_id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.CreateLink(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "job_id", jobID)
}

func TestCreateEnvironmentUnlink(t *testing.T) {
	req := models.DeleteEnvironmentLinkRequest{
		SourceEnvironmentID: "env_id1",
		DestEnvironmentID:   "env_id2",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/environment/env_id1/link/env_id2")

		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "job_id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.DeleteLink(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "job_id", jobID)
}
