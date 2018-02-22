package client

import (
	"net/http"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestCreateLoadBalancer(t *testing.T) {
	healthCheck := models.HealthCheck{
		Target:             "TCP:80",
		Interval:           30,
		Timeout:            5,
		HealthyThreshold:   10,
		UnhealthyThreshold: 2,
	}

	ports := []models.Port{
		{
			HostPort:        443,
			ContainerPort:   80,
			Protocol:        "https",
			CertificateName: "cert_name",
		},
		{
			HostPort:        8000,
			ContainerPort:   8000,
			Protocol:        "http",
			CertificateName: "",
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "POST")
		testutils.AssertEqual(t, r.URL.Path, "/loadbalancer/")

		var req models.CreateLoadBalancerRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.LoadBalancerName, "name")
		testutils.AssertEqual(t, req.EnvironmentID, "environmentID")
		testutils.AssertEqual(t, req.IsPublic, true)
		testutils.AssertEqual(t, req.HealthCheck, healthCheck)
		testutils.AssertEqual(t, req.Ports, ports)
		testutils.AssertEqual(t, req.IdleTimeout, 60)

		MarshalAndWrite(t, w, models.LoadBalancer{LoadBalancerID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	loadBalancer, err := client.CreateLoadBalancer("name", "environmentID", healthCheck, ports, true, 60)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, loadBalancer.LoadBalancerID, "id")
}

func TestDeleteLoadBalancer(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "DELETE")
		testutils.AssertEqual(t, r.URL.Path, "/loadbalancer/id")

		headers := map[string]string{
			"Location": "/job/jobid",
			"X-JobID":  "jobid",
		}

		MarshalAndWriteHeader(t, w, "", headers, 202)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.DeleteLoadBalancer("id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, jobID, "jobid")
}

func TestGetLoadBalancer(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/loadbalancer/id")

		MarshalAndWrite(t, w, models.LoadBalancer{LoadBalancerID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	loadBalancer, err := client.GetLoadBalancer("id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, loadBalancer.LoadBalancerID, "id")
}

func TestListLoadBalancers(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/loadbalancer/")

		loadBalancers := []models.LoadBalancerSummary{
			{LoadBalancerID: "id1"},
			{LoadBalancerID: "id2"},
		}

		MarshalAndWrite(t, w, loadBalancers, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	loadBalancers, err := client.ListLoadBalancers()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(loadBalancers), 2)
	testutils.AssertEqual(t, loadBalancers[0].LoadBalancerID, "id1")
	testutils.AssertEqual(t, loadBalancers[1].LoadBalancerID, "id2")
}

func TestUpdateLoadBalancerHealthCheck(t *testing.T) {
	healthCheck := models.HealthCheck{
		Target:             "TCP:80",
		Interval:           30,
		Timeout:            5,
		HealthyThreshold:   10,
		UnhealthyThreshold: 2,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "PUT")
		testutils.AssertEqual(t, r.URL.Path, "/loadbalancer/id/healthcheck")

		var req models.UpdateLoadBalancerHealthCheckRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.HealthCheck, healthCheck)

		MarshalAndWrite(t, w, models.LoadBalancer{LoadBalancerID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	loadBalancer, err := client.UpdateLoadBalancerHealthCheck("id", healthCheck)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, loadBalancer.LoadBalancerID, "id")
}

func TestUpdateLoadBalancerPorts(t *testing.T) {
	ports := []models.Port{
		{
			HostPort:        443,
			ContainerPort:   80,
			Protocol:        "https",
			CertificateName: "cert_name",
		},
		{
			HostPort:        8000,
			ContainerPort:   8000,
			Protocol:        "http",
			CertificateName: "",
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "PUT")
		testutils.AssertEqual(t, r.URL.Path, "/loadbalancer/id/ports")

		var req models.UpdateLoadBalancerPortsRequest
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.Ports, ports)

		MarshalAndWrite(t, w, models.LoadBalancer{LoadBalancerID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	loadBalancer, err := client.UpdateLoadBalancerPorts("id", ports)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, loadBalancer.LoadBalancerID, "id")
}
