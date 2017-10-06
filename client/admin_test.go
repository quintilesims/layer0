package client

import (
	"net/http"
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
