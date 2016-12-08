package client

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"net/http"
	"testing"
	"time"
)

func TestCreateService(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "POST")
		testutils.AssertEqual(t, r.URL.Path, "/service/")

		var req models.CreateServiceRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.ServiceName, "name")
		testutils.AssertEqual(t, req.EnvironmentID, "environmentID")
		testutils.AssertEqual(t, req.DeployID, "deployID")
		testutils.AssertEqual(t, req.LoadBalancerID, "loadBalancerID")

		MarshalAndWrite(t, w, models.Service{ServiceID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	service, err := client.CreateService("name", "environmentID", "deployID", "loadBalancerID")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, service.ServiceID, "id")
}

func TestDeleteService(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "DELETE")
		testutils.AssertEqual(t, r.URL.Path, "/service/id")

		headers := map[string]string{
			"Location": "/job/jobid",
			"X-JobID":  "jobid",
		}

		MarshalAndWriteHeader(t, w, "", headers, 202)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.DeleteService("id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, jobID, "jobid")
}

func TestGetService(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/service/id")

		MarshalAndWrite(t, w, models.Service{ServiceID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	service, err := client.GetService("id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, service.ServiceID, "id")
}

func TestGetServiceLogs(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/service/id/logs")
		testutils.AssertEqual(t, r.URL.RawQuery, "tail=100")

		logs := []models.LogFile{
			{Name: "name1"},
			{Name: "name2"},
		}

		MarshalAndWrite(t, w, logs, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	logs, err := client.GetServiceLogs("id", 100)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(logs), 2)
	testutils.AssertEqual(t, logs[0].Name, "name1")
	testutils.AssertEqual(t, logs[1].Name, "name2")
}

func TestListServices(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/service/")

		services := []models.Service{
			{ServiceID: "id1"},
			{ServiceID: "id2"},
		}

		MarshalAndWrite(t, w, services, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	services, err := client.ListServices()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(services), 2)
	testutils.AssertEqual(t, services[0].ServiceID, "id1")
	testutils.AssertEqual(t, services[1].ServiceID, "id2")
}

func TestScaleService(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "PUT")
		testutils.AssertEqual(t, r.URL.Path, "/service/id/scale")

		var req models.ScaleServiceRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.DesiredCount, int64(2))

		MarshalAndWrite(t, w, models.Service{ServiceID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	service, err := client.ScaleService("id", 2)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, service.ServiceID, "id")
}

func TestUpdateService(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "PUT")
		testutils.AssertEqual(t, r.URL.Path, "/service/id/update")

		var req models.UpdateServiceRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.DeployID, "deployID")

		MarshalAndWrite(t, w, models.Service{ServiceID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	service, err := client.UpdateService("id", "deployID")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, service.ServiceID, "id")
}

func TestWaitForDeployment(t *testing.T) {
	var count int

	handler := func(w http.ResponseWriter, r *http.Request) {
		count++
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/service/id")

		var runningCount int64 = 1
		var desiredCount int64 = 2

		// simulate flapping success
		if count == 0 || count > 3 {
			runningCount = 2
		}

		service := models.Service{
			ServiceID: "id",
			Deployments: []models.Deployment{
				{DesiredCount: desiredCount, RunningCount: runningCount},
				{DesiredCount: desiredCount, RunningCount: runningCount},
			},
		}

		MarshalAndWrite(t, w, service, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	service, err := client.WaitForDeployment("id", time.Minute*15)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, service.ServiceID, "id")
	if count < REQUIRED_SUCCESS_WAIT_COUNT {
		t.Fatalf("Retry count was less than required (%d)", count)
	}
}

func TestWaitForDeployment_timeout(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		service := models.Service{
			Deployments: []models.Deployment{
				{DesiredCount: 1, RunningCount: 0},
			},
		}

		MarshalAndWrite(t, w, service, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if _, err := client.WaitForDeployment("id", time.Millisecond); err == nil {
		t.Fatal("Error was nil!")
	}
}
