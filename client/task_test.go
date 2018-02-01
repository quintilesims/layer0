package client

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	req := models.CreateTaskRequest{
		TaskName:      "tsk_name",
		EnvironmentID: "env_id",
		DeployID:      "dpl_id",
		ContainerOverrides: []models.ContainerOverride{
			{ContainerName: "c1", EnvironmentOverrides: map[string]string{"k1": "v1"}},
			{ContainerName: "c2", EnvironmentOverrides: map[string]string{"k2": "v2"}},
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/task")

		var body models.CreateTaskRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateEntityResponse{EntityID: "tsk_id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	taskID, err := client.CreateTask(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "tsk_id", taskID)
}

func TestDeleteTask(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/task/tsk_id")

		MarshalAndWrite(t, w, nil, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteTask("tsk_id"); err != nil {
		t.Fatal(err)
	}
}

func TestListTasks(t *testing.T) {
	expected := []models.TaskSummary{
		{
			TaskID:          "tsk_id1",
			TaskName:        "tsk_name1",
			EnvironmentID:   "env_id1",
			EnvironmentName: "env_name1",
		},
		{
			TaskID:          "tsk_id2",
			TaskName:        "tskd_name2",
			EnvironmentID:   "env_id2",
			EnvironmentName: "env_name2",
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/task")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ListTasks()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestReadTask(t *testing.T) {
	expected := models.Task{
		TaskID:          "tsk_id",
		TaskName:        "tsk_name",
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		DeployID:        "dpl_id",
		DeployName:      "dpl_name",
		DeployVersion:   "1",
		Status:          "RUNNING",
		Containers: []models.Container{
			{ContainerName: "c1", Status: "RUNNING", ExitCode: 0},
			{ContainerName: "c2", Status: "STOPPED", ExitCode: 1},
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/task/tsk_id")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadTask("tsk_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, *result)
}

func TestReadTaskLogs(t *testing.T) {
	expected := []models.LogFile{
		{
			ContainerName: "apline",
			Lines:         []string{"hello", "world"},
		},
	}

	query := url.Values{}
	query.Set(models.LogQueryParamTail, "100")
	query.Set(models.LogQueryParamStart, "2000-01-01 00:00")
	query.Set(models.LogQueryParamEnd, "2000-01-01 12:12")

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/task/tsk_id/logs")
		assert.Equal(t, query, r.URL.Query())

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadTaskLogs("tsk_id", query)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}
