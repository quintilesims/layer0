package client

import (
	"net/http"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDeleteJob(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/job/jid")

		MarshalAndWrite(t, w, nil, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteJob("jid"); err != nil {
		t.Fatal(err)
	}
}

func TestListJobs(t *testing.T) {
	expected := []*models.Job{
		{JobID: "jid1"},
		{JobID: "jid2"},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/job")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ListJobs()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestReadJob(t *testing.T) {
	expected := &models.Job{
		JobID:  "jid",
		Type:   "CreateEnvironmentJob",
		Status: "InProgress",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/job/jid")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadJob("jid")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}
