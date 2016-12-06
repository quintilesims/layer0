package client

import (
	// "encoding/base64"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/common/types"
	"net/http"
	"testing"
)

func TestDeleteJob(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "DELETE")
		testutils.AssertEqual(t, r.URL.Path, "/job/id")

		MarshalAndWrite(t, w, "", 202)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteJob("id"); err != nil {
		t.Fatal(err)
	}
}

func TestGetJob(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/job/id")

		MarshalAndWrite(t, w, models.Job{JobID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	job, err := client.GetJob("id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, job.JobID, "id")
}

func TestListJobs(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/job/")

		jobs := []models.Job{
			{JobID: "id1"},
			{JobID: "id2"},
		}

		MarshalAndWrite(t, w, jobs, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobs, err := client.ListJobs()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(jobs), 2)
	testutils.AssertEqual(t, jobs[0].JobID, "id1")
	testutils.AssertEqual(t, jobs[1].JobID, "id2")
}

func TestWaitForJob(t *testing.T) {
	count := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/job/id")

		jobStatus := int64(types.InProgress)
		if count > 0 {
			jobStatus = int64(types.Completed)
		}

		job := models.Job{JobID: "id", JobStatus: jobStatus}

		MarshalAndWrite(t, w, job, 200)
		count++
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.WaitForJob("id", 0); err != nil {
		t.Fatal(err)
	}
}

func TestWaitForJobError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/job/id")

		jobStatus := int64(types.Error)

		job := models.Job{JobID: "id", JobStatus: jobStatus}

		MarshalAndWrite(t, w, job, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.WaitForJob("id", 0); err == nil {
		t.Fatalf("Error was nil!")
	}
}
