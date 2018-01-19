package client

import (
	"net/http"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeploy(t *testing.T) {
	req := models.CreateDeployRequest{
		DeployName: "dpl_name",
		DeployFile: []byte("content"),
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/deploy")

		var body models.CreateDeployRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateEntityResponse{EntityID: "dpl_id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	deployID, err := client.CreateDeploy(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "dpl_id", deployID)
}

func TestDeleteDeploy(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/deploy/dpl_id")

		MarshalAndWrite(t, w, nil, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteDeploy("dpl_id"); err != nil {
		t.Fatal(err)
	}
}

func TestListDeploys(t *testing.T) {
	expected := []models.DeploySummary{
		{
			DeployID:   "dpl_id1",
			DeployName: "dpl_name1",
			Version:    "1",
		},
		{
			DeployID:   "dpl_id2",
			DeployName: "dpld_name2",
			Version:    "2",
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/deploy")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ListDeploys()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestReadDeploy(t *testing.T) {
	expected := models.Deploy{
		DeployID:   "dpl_id",
		DeployName: "dpl_name",
		Version:    "1",
		DeployFile: []byte("content"),
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/deploy/dpl_id")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadDeploy("dpl_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}
