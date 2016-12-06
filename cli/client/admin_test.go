package client

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"net/http"
	"testing"
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
