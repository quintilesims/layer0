package rclient

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestBuildJSONRequest(t *testing.T) {
	req, err := BuildJSONRequest("GET", "www.domain.com/path", "body", Header("name", "val"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "www.domain.com/path", req.URL.String())
	assert.Equal(t, "application/json", req.Header.Get("content-type"))
	assert.Equal(t, "val", req.Header.Get("name"))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "\"body\"\n", string(body))
}

func TestBuildJSONRequest_optionError(t *testing.T) {
	option := func(req *http.Request) error {
		return errors.New("some error")
	}

	if _, err := BuildJSONRequest("GET", "www.domain.com/path", "body", option); err == nil {
		t.Fatal("Error was nil!")
	}
}
