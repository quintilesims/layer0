package client

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	expected := &models.APIConfig{
		Instance:       "instance",
		VPCID:          "vpcid",
		Version:        "version",
		PublicSubnets:  []string{},
		PrivateSubnets: []string{},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/admin/config")

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadConfig()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestReadAdminLogs(t *testing.T) {
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
		assert.Equal(t, r.URL.Path, "/admin/logs")
		assert.Equal(t, query, r.URL.Query())

		MarshalAndWrite(t, w, expected, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	result, err := client.ReadAdminLogs(query)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}
