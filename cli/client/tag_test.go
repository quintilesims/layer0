package client

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"net/http"
	"testing"
)

func TestGetTags(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testutils.AssertEqual(t, r.Method, "GET")
		testutils.AssertEqual(t, r.URL.Path, "/tag")

		query := r.URL.Query()
		testutils.AssertEqual(t, query.Get("type"), "some_type")
		testutils.AssertEqual(t, query.Get("fuzz"), "some_fuzz")
		testutils.AssertEqual(t, query.Get("version"), "some_version")
		testutils.AssertEqual(t, query.Get("key"), "val")

		tags := []models.EntityWithTags{
			{EntityID: "id1"},
			{EntityID: "id2"},
		}

		MarshalAndWrite(t, w, tags, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	params := map[string]string{
		"type":    "some_type",
		"fuzz":    "some_fuzz",
		"version": "some_version",
		"key":     "val",
	}

	tags, err := client.GetTags(params)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(tags), 2)
	testutils.AssertEqual(t, tags[0].EntityID, "id1")
	testutils.AssertEqual(t, tags[1].EntityID, "id2")
}
