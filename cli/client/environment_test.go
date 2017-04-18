package client

import (
	"net/http"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestCreateEnvironment(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "POST")
		testutils.AssertEqual(t, r.URL.Path, "/environment/")

		var req models.CreateEnvironmentRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.EnvironmentName, "name")
		testutils.AssertEqual(t, req.InstanceSize, "m3.medium")
		testutils.AssertEqual(t, req.MinClusterCount, 2)
		testutils.AssertEqual(t, req.UserDataTemplate, []byte("user_data"))
		testutils.AssertEqual(t, req.OperatingSystem, "linux")
		testutils.AssertEqual(t, req.AMIID, "ami")

		MarshalAndWrite(t, w, models.Environment{EnvironmentID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	environment, err := client.CreateEnvironment("name", "m3.medium", 2, []byte("user_data"), "linux", "ami")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, environment.EnvironmentID, "id")
}

func TestDeleteEnvironment(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "DELETE")
		testutils.AssertEqual(t, r.URL.Path, "/environment/id")

		headers := map[string]string{
			"Location": "/job/jobid",
			"X-JobID":  "jobid",
		}

		MarshalAndWriteHeader(t, w, "", headers, 202)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.DeleteEnvironment("id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, jobID, "jobid")
}

func TestGetEnvironment(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/environment/id")

		MarshalAndWrite(t, w, models.Environment{EnvironmentID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	environment, err := client.GetEnvironment("id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, environment.EnvironmentID, "id")
}

func TestListEnvironments(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/environment/")

		environments := []models.EnvironmentSummary{
			{EnvironmentID: "id1"},
			{EnvironmentID: "id2"},
		}

		MarshalAndWrite(t, w, environments, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	environments, err := client.ListEnvironments()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(environments), 2)
	testutils.AssertEqual(t, environments[0].EnvironmentID, "id1")
	testutils.AssertEqual(t, environments[1].EnvironmentID, "id2")
}

func TestUpdateEnvironment(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "PUT")
		testutils.AssertEqual(t, r.URL.Path, "/environment/id")

		var req models.UpdateEnvironmentRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.MinClusterCount, 2)

		MarshalAndWrite(t, w, models.Environment{EnvironmentID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	environment, err := client.UpdateEnvironment("id", 2)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, environment.EnvironmentID, "id")
}

func TestCreateLink(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "POST")
		testutils.AssertEqual(t, r.URL.Path, "/environment/id1/link")

		var req models.CreateEnvironmentLinkRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.EnvironmentID, "id2")

		MarshalAndWrite(t, w, "", 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.CreateLink("id1", "id2"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateUnink(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "DELETE")
		testutils.AssertEqual(t, r.URL.Path, "/environment/id1/link/id2")

		MarshalAndWrite(t, w, "", 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteLink("id1", "id2"); err != nil {
		t.Fatal(err)
	}
}