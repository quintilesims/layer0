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
		InstanceType:     "t2.small",
		MinScale:         1,
		MaxScale:         5,
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
		MarshalAndWrite(t, w, models.CreateEntityResponse{EntityID: "env_id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	environmentID, err := client.CreateEnvironment(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "env_id", environmentID)
}

func TestDeleteEnvironment(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/environment/env_id")

		MarshalAndWrite(t, w, models.CreateEntityResponse{EntityID: "env_id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteEnvironment("env_id"); err != nil {
		t.Fatal(err)
	}
}

func TestListEnvironments(t *testing.T) {
	expected := []*models.EnvironmentSummary{
		{EnvironmentID: "env_id1"},
		{EnvironmentID: "env_id2"},
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
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/environment/env_id")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadEnvironment("env_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestUpdateEnvironment(t *testing.T) {
	minScale := 1
	maxScale := 2
	links := []string{"env_id2", "env_id3"}

	req := models.UpdateEnvironmentRequest{
		MinScale: &minScale,
		MaxScale: &maxScale,
		Links:    &links,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PATCH")
		assert.Equal(t, r.URL.Path, "/environment/env_id1")

		var body models.UpdateEnvironmentRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.UpdateEnvironment("env_id", req); err != nil {
		t.Fatal(err)
	}
}
