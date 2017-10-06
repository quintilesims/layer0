package client

import (
	"net/http"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateLoadBalancer(t *testing.T) {
	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "name",
		EnvironmentID:    "eid",
		IsPublic:         true,
		Ports:            []models.Port{},
		HealthCheck:      models.HealthCheck{},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/loadbalancer")

		var body models.CreateLoadBalancerRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "jid"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.CreateLoadBalancer(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "jid", jobID)
}

func TestDeleteLoadBalancer(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/loadbalancer/lid")

		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "jid"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.DeleteLoadBalancer("lid")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, jobID, "jid")
}

func TestListLoadBalancers(t *testing.T) {
	expected := []*models.LoadBalancerSummary{
		{LoadBalancerID: "lid1"},
		{LoadBalancerID: "lid2"},
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
	expected := &models.LoadBalancer{
		LoadBalancerID:   "lid",
		LoadBalancerName: "lname",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/loadbalancer/lid")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadLoadBalancer("lid")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestUpdateLoadBalancer(t *testing.T) {
	ports := []models.Port{}
	healthCheck := models.HealthCheck{}
	req := models.UpdateLoadBalancerRequest{
		LoadBalancerID: "lid",
		Ports:          &ports,
		HealthCheck:    &healthCheck,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/loadbalancer")

		var body models.UpdateLoadBalancerRequest
		Unmarshal(t, r, &body)

		assert.Equal(t, req, body)
		MarshalAndWrite(t, w, models.CreateJobResponse{JobID: "jid"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.UpdateLoadBalancer(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "jid", jobID)
}
