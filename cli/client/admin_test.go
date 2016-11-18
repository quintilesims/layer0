package client

import (
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/testutils"
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
