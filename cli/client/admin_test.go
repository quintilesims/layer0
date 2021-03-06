package client

import (
	"net/http"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestGetVersion(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/admin/version")

		MarshalAndWrite(t, w, "v1.2.3", 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	version, err := client.GetVersion()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, version, "v1.2.3")
}

func TestGetConfig(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/admin/config")

		MarshalAndWrite(t, w, models.APIConfig{VPCID: "vpc"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	config, err := client.GetConfig()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, config.VPCID, "vpc")
}

func TestUpdateSQL(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "POST")
		testutils.AssertEqual(t, r.URL.Path, "/admin/sql")

		var req models.SQLVersion
		Unmarshal(t, r, &req)

		testutils.AssertEqual(t, req.Version, "latest")

		MarshalAndWrite(t, w, "", 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.UpdateSQL(); err != nil {
		t.Fatal(err)
	}
}

func TestRunScaler(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "PUT")
		testutils.AssertEqual(t, r.URL.Path, "/admin/scale/id")

		MarshalAndWrite(t, w, models.ScalerRunInfo{EnvironmentID: "id"}, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	output, err := client.RunScaler("id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, output.EnvironmentID, "id")
}
