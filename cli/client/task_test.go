package client

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"net/http"
	"testing"
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
		testutils.AssertEqual(t, req.Copies, int64(2))
		testutils.AssertEqual(t, req.ContainerOverrides, overrides)

		MarshalAndWrite(t, w, models.Task{TaskID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	task, err := client.CreateTask("name", "environmentID", "deployID", 2, overrides)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, task.TaskID, "id")
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
		testutils.AssertEqual(t, r.URL.RawQuery, "tail=100")

		logs := []models.LogFile{
			{Name: "name1"},
			{Name: "name2"},
		}

		MarshalAndWrite(t, w, logs, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	logs, err := client.GetTaskLogs("id", 100)
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

		tasks := []models.Task{
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
