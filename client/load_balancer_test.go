package client

import (
	"net/http"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateLoadBalancer(t *testing.T) {
	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		IsPublic:         true,
		Ports: []models.Port{
			{HostPort: 443, ContainerPort: 80, Protocol: "https", CertificateName: "cert"},
			{HostPort: 22, ContainerPort: 22, Protocol: "tcp"},
		},
		HealthCheck: models.HealthCheck{
			Target:             "tcp:80",
			Interval:           1,
			Timeout:            2,
			HealthyThreshold:   3,
			UnhealthyThreshold: 4,
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/loadbalancer")

		var body models.CreateLoadBalancerRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateEntityResponse{EntityID: "lb_id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	loadBalancerID, err := client.CreateLoadBalancer(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "lb_id", loadBalancerID)
}

func TestDeleteLoadBalancer(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/loadbalancer/lb_id")

		MarshalAndWrite(t, w, nil, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteLoadBalancer("lb_id"); err != nil {
		t.Fatal(err)
	}
}

func TestListLoadBalancers(t *testing.T) {
	expected := []models.LoadBalancerSummary{
		{
			LoadBalancerID:   "lb_id1",
			LoadBalancerName: "lb_name1",
			EnvironmentID:    "env_id1",
			EnvironmentName:  "env_name1",
		},
		{
			LoadBalancerID:   "lb_id2",
			LoadBalancerName: "lbd_name2",
			EnvironmentID:    "env_id2",
			EnvironmentName:  "env_name2",
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/loadbalancer")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ListLoadBalancers()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestReadLoadBalancer(t *testing.T) {
	expected := models.LoadBalancer{
		LoadBalancerID:   "lb_id",
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		EnvironmentName:  "env_name",
		ServiceID:        "svc_id",
		ServiceName:      "svc_name",
		IsPublic:         true,
		URL:              "url",
		Ports: []models.Port{
			{HostPort: 443, ContainerPort: 80, Protocol: "https", CertificateName: "cert"},
			{HostPort: 22, ContainerPort: 22, Protocol: "tcp"},
		},
		HealthCheck: models.HealthCheck{
			Target:             "tcp:80",
			Interval:           1,
			Timeout:            2,
			HealthyThreshold:   3,
			UnhealthyThreshold: 4,
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/loadbalancer/lb_id")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadLoadBalancer("lb_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, *result)
}

func TestUpdateLoadBalancer(t *testing.T) {
	ports := []models.Port{
		{HostPort: 443, ContainerPort: 80, Protocol: "https", CertificateName: "cert"},
		{HostPort: 22, ContainerPort: 22, Protocol: "tcp"},
	}

	healthCheck := models.HealthCheck{
		Target:             "tcp:80",
		Interval:           1,
		Timeout:            2,
		HealthyThreshold:   3,
		UnhealthyThreshold: 4,
	}

	req := models.UpdateLoadBalancerRequest{
		Ports:       &ports,
		HealthCheck: &healthCheck,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PATCH")
		assert.Equal(t, r.URL.Path, "/loadbalancer/lb_id")

		var body models.UpdateLoadBalancerRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, nil, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.UpdateLoadBalancer("lb_id", req); err != nil {
		t.Fatal(err)
	}
}
