package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type Handler func(w http.ResponseWriter, r *http.Request)

func newClientAndServer(handler Handler) (*APIClient, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(handler))
	client := NewAPIClient(Config{
		Endpoint: server.URL,
	})

	return client, server
}

func MarshalAndWrite(t *testing.T, w http.ResponseWriter, body interface{}, status int) {
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
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
