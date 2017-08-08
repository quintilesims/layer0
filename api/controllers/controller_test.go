package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	restful "github.com/emicklei/go-restful"
)

func newRequest(t *testing.T, body interface{}, params map[string]string) *restful.Request {
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("", "", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	restfulRequest := restful.NewRequest(req)
	for key, val := range params {
		restfulRequest.PathParameters()[key] = val
	}

	req.Header.Set("Content-Type", "application/json")
	return restfulRequest
}

func unmarshalBody(t *testing.T, b []byte, v interface{}) {
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatal(err)
	}
}
