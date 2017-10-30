package rclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type person struct {
	Name string
	Age  int
}

func readPerson(t *testing.T, r *http.Request) person {
	var p person
	read(t, r, &p)
	return p
}

func newClientAndServer(t *testing.T, handler http.HandlerFunc, options ...ClientOption) (*RestClient, *httptest.Server) {
	server := httptest.NewServer(handler)
	client := NewRestClient(server.URL, options...)

	return client, server
}

func write(t *testing.T, w http.ResponseWriter, status int, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	w.WriteHeader(status)
	fmt.Fprintln(w, string(b))
}

func read(t *testing.T, r *http.Request, v interface{}) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		t.Fatal(err)
	}
}
