package client

import (
	"encoding/json"
	"fmt"
	"github.com/quintilesims/layer0/common/testutils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Handler func(w http.ResponseWriter, r *http.Request)

func newClientAndServer(handler Handler) (*APIClient, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(handler))
	client := NewAPIClient(Config{
		Endpoint: server.URL,
		Clock:    &testutils.StubClock{},
	})

	return client, server
}

func MarshalAndWrite(t *testing.T, w http.ResponseWriter, body interface{}, status int) {
	MarshalAndWriteHeader(t, w, body, nil, status)
}

func MarshalAndWriteHeader(t *testing.T, w http.ResponseWriter, body interface{}, headers map[string]string, status int) {
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	for key, val := range headers {
		w.Header().Set(key, val)
	}

	w.WriteHeader(status)
	fmt.Fprintln(w, string(b))
}

func Unmarshal(t *testing.T, r *http.Request, content interface{}) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(body, &content); err != nil {
		t.Fatal(err)
	}
}
