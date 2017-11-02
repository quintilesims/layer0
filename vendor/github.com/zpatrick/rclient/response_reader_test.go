package rclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadJSONResponse(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("\"body\"")),
	}

	var v string
	if err := ReadJSONResponse(resp, &v); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "body", v)
}

func TestReadJSONResponseNilV(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("\"body\"")),
	}

	if err := ReadJSONResponse(resp, nil); err != nil {
		t.Fatal(err)
	}
}

func TestReadJSONResponseError_statusCode(t *testing.T) {
	codes := []int{0, 199, 300, 399, 400, 499, 500, 599}

	for _, c := range codes {
		resp := &http.Response{
			StatusCode: c,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}

		if err := ReadJSONResponse(resp, nil); err == nil {
			t.Fatalf("%d: Error was nil!", c)
		}
	}
}

func TestReadJSONResponseError_invalidJSON(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("some_invalid_json")),
	}

	var p person
	if err := ReadJSONResponse(resp, &p); err == nil {
		t.Fatalf("Error was nil!")
	}
}
