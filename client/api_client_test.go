package client

import (
	"net/http"
	"testing"

	"github.com/quintilesims/layer0/common/errors"
	"github.com/stretchr/testify/assert"
)

func TestAPIClientRetry(t *testing.T) {
	first := true
	handler := func(w http.ResponseWriter, r *http.Request) {
		if first {
			first = false
			serverError := errors.Newf(errors.EventualConsistencyError, "")
			MarshalAndWrite(t, w, serverError.Model(), 500)
			return
		}

		MarshalAndWrite(t, w, nil, 200)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteEnvironment(""); err != nil {
		t.Fatal(err)
	}
}

func TestAPIClientRetryError(t *testing.T) {
	var calls int
	handler := func(w http.ResponseWriter, r *http.Request) {
		calls++
		serverError := errors.Newf(errors.EventualConsistencyError, "")
		MarshalAndWrite(t, w, serverError.Model(), 500)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	if err := client.DeleteEnvironment(""); err == nil {
		t.Fatal("Error was nil!")
	}

	assert.Equal(t, MaxRetries, calls)
}
