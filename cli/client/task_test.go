package client

import (
	"net/http"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestCreateTask(t *testing.T) {
	overrides := []models.ContainerOverride{{
		ContainerName:        "container",
		EnvironmentOverrides: map[string]string{"key": "val"},
	}}

	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "POST")
		testutils.AssertEqual(t, r.URL.Path, "/task/")

		var req models.CreateTaskRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.TaskName, "name")
		testutils.AssertEqual(t, req.EnvironmentID, "environmentID")
		testutils.AssertEqual(t, req.DeployID, "deployID")
		testutils.AssertEqual(t, req.ContainerOverrides, overrides)

		headers := map[string]string{
			"Location": "/job/jobid",
			"X-JobID":  "jobid",
		}

		MarshalAndWriteHeader(t, w, "", headers, 202)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.CreateTask("name", "environmentID", "deployID", overrides)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, jobID, "jobid")
}

func TestDeleteTask(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "DELETE")
		testutils.AssertEqual(t, r.URL.Path, "/task/id")

		MarshalAndWrite(t, w, "", 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteTask("id"); err != nil {
		t.Fatal(err)
	}
}

func TestGetTask(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/task/id")

		MarshalAndWrite(t, w, models.Task{TaskID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	task, err := client.GetTask("id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, task.TaskID, "id")
}

func TestGetTaskLogs(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/task/id/logs")
		testutils.AssertEqual(t, r.URL.Query().Get("tail"), "100")
		testutils.AssertEqual(t, r.URL.Query().Get("start"), "2001-01-01 01:01")
		testutils.AssertEqual(t, r.URL.Query().Get("end"), "2012-12-12 12:12")

		logs := []models.LogFile{
			{Name: "name1"},
			{Name: "name2"},
		}

		MarshalAndWrite(t, w, logs, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	logs, err := client.GetTaskLogs("id", "2001-01-01 01:01", "2012-12-12 12:12", 100)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(logs), 2)
	testutils.AssertEqual(t, logs[0].Name, "name1")
	testutils.AssertEqual(t, logs[1].Name, "name2")
}

func TestListTasks(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/task/")

		tasks := []models.TaskSummary{
			{TaskID: "id1"},
			{TaskID: "id2"},
		}

		MarshalAndWrite(t, w, tasks, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	tasks, err := client.ListTasks()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(tasks), 2)
	testutils.AssertEqual(t, tasks[0].TaskID, "id1")
	testutils.AssertEqual(t, tasks[1].TaskID, "id2")
}
