package client

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"net/http"
	"testing"
)

func TestCreateDeploy(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "POST")
		testutils.AssertEqual(t, r.URL.Path, "/deploy")

		var req models.CreateDeployRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.DeployName, "name")
		testutils.AssertEqual(t, req.Dockerrun, []byte("dockerrun"))

		MarshalAndWrite(t, w, models.Deploy{DeployID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	deploy, err := client.CreateDeploy("name", []byte("dockerrun"))
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, deploy.DeployID, "id")
}

func TestDeleteDeploy(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "DELETE")
		testutils.AssertEqual(t, r.URL.Path, "/deploy/id")

		MarshalAndWrite(t, w, "", 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteDeploy("id"); err != nil {
		t.Fatal(err)
	}
}

func TestGetDeploy(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/deploy/id")

		MarshalAndWrite(t, w, models.Deploy{DeployID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	deploy, err := client.GetDeploy("id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, deploy.DeployID, "id")
}

func TestListDeploys(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/deploy/")

		deploys := []models.Deploy{
			{DeployID: "id1"},
			{DeployID: "id2"},
		}

		MarshalAndWrite(t, w, deploys, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	deploys, err := client.ListDeploys()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(deploys), 2)
	testutils.AssertEqual(t, deploys[0].DeployID, "id1")
	testutils.AssertEqual(t, deploys[1].DeployID, "id2")
}
