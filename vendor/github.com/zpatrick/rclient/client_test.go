package rclient

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var JohnDoe = person{Name: "John Doe", Age: 35}

func TestClientDelete(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/people/john", r.URL.Path)

		write(t, w, 200, nil)
	}

	client, server := newClientAndServer(t, handler)
	defer server.Close()

	if err := client.Delete("/people/john", nil, nil); err != nil {
		t.Error(err)
	}
}

func TestClientDeleteWithBody(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/people", r.URL.Path)
		assert.Equal(t, JohnDoe, readPerson(t, r))

		write(t, w, 200, JohnDoe)
	}

	client, server := newClientAndServer(t, handler)
	defer server.Close()

	var p person
	if err := client.Delete("/people", JohnDoe, &p); err != nil {
		t.Error(err)
	}

	assert.Equal(t, JohnDoe, p)
}

func TestClientGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/people/john", r.URL.Path)

		write(t, w, 200, JohnDoe)
	}

	client, server := newClientAndServer(t, handler)
	defer server.Close()

	var p person
	if err := client.Get("/people/john", &p); err != nil {
		t.Error(err)
	}

	assert.Equal(t, JohnDoe, p)
}

func TestClientPatch(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/people/john", r.URL.Path)
		assert.Equal(t, JohnDoe, readPerson(t, r))

		write(t, w, 200, JohnDoe)
	}

	client, server := newClientAndServer(t, handler)
	defer server.Close()

	var p person
	if err := client.Patch("/people/john", JohnDoe, &p); err != nil {
		t.Error(err)
	}

	assert.Equal(t, JohnDoe, p)
}

func TestClientPost(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/people", r.URL.Path)
		assert.Equal(t, JohnDoe, readPerson(t, r))

		write(t, w, 201, JohnDoe)
	}

	client, server := newClientAndServer(t, handler)
	defer server.Close()

	var p person
	if err := client.Post("/people", JohnDoe, &p); err != nil {
		t.Error(err)
	}

	assert.Equal(t, JohnDoe, p)
}

func TestClientPut(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/people/john", r.URL.Path)
		assert.Equal(t, JohnDoe, readPerson(t, r))

		write(t, w, 200, JohnDoe)
	}

	client, server := newClientAndServer(t, handler)
	defer server.Close()

	var p person
	if err := client.Put("/people/john", JohnDoe, &p); err != nil {
		t.Error(err)
	}

	assert.Equal(t, JohnDoe, p)
}

func TestClientDo(t *testing.T) {
	builder := func(method, url string, body interface{}, options ...RequestOption) (*http.Request, error) {
		assert.Equal(t, "POST", method)
		assert.Equal(t, "https://domain.com/path", url)
		assert.Equal(t, "body", body)
		assert.Len(t, options, 0)

		return nil, nil
	}

	doer := RequestDoerFunc(func(*http.Request) (*http.Response, error) {
		return nil, nil
	})

	var p person
	reader := func(resp *http.Response, v interface{}) error {
		assert.Equal(t, p, v)
		return nil
	}

	client := NewRestClient("https://domain.com", Builder(builder), Doer(doer), Reader(reader))
	if err := client.Post("/path", "body", p); err != nil {
		t.Fatal(err)
	}
}

func TestClientBuilderError(t *testing.T) {
	builder := func(string, string, interface{}, ...RequestOption) (*http.Request, error) {
		return nil, errors.New("some error")
	}

	client := NewRestClient("", Builder(builder))
	if err := client.Get("/path", nil); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestClientDoerError(t *testing.T) {
	builder := func(string, string, interface{}, ...RequestOption) (*http.Request, error) {
		return nil, nil
	}

	doer := RequestDoerFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("some error")
	})

	client := NewRestClient("", Builder(builder), Doer(doer))
	if err := client.Get("/path", nil); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestClientReaderError(t *testing.T) {
	builder := func(string, string, interface{}, ...RequestOption) (*http.Request, error) {
		return nil, nil
	}

	doer := RequestDoerFunc(func(*http.Request) (*http.Response, error) {
		return nil, nil
	})

	reader := func(*http.Response, interface{}) error {
		return errors.New("some error")
	}

	client := NewRestClient("", Builder(builder), Doer(doer), Reader(reader))
	if err := client.Get("/path", nil); err == nil {
		t.Fatal("Error was nil!")
	}
}
