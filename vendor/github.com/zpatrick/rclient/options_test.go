package rclient

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	req, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := BasicAuth("user", "pass")(req); err != nil {
		t.Fatal(err)
	}

	// user:pass base64 encoded
	assert.Equal(t, "Basic dXNlcjpwYXNz", req.Header.Get("Authorization"))
}

func TestHeader(t *testing.T) {
	req, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := Header("name", "val")(req); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "val", req.Header.Get("name"))
}

func TestHeaders(t *testing.T) {
	req, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := map[string]string{"name1": "v1", "name2": "v2"}
	if err := Headers(h)(req); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "v1", req.Header.Get("name1"))
	assert.Equal(t, "v2", req.Header.Get("name2"))
}

func TestQuery(t *testing.T) {
	req, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := url.Values{}
	q.Set("k1", "v1")
	q.Set("k2", "v2")

	if err := Query(q)(req); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "k1=v1&k2=v2", req.URL.RawQuery)
}
