package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zpatrick/fireball"
)

func newFireballContext(t *testing.T, body interface{}, params map[string]string) *fireball.Context {
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("", "", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	return &fireball.Context{
		Request:       req,
		PathVariables: params,
	}
}

func unmarshalBody(t *testing.T, resp fireball.Response, v interface{}) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	resp.Write(recorder, nil)

	if v != nil {
		if err := json.Unmarshal(recorder.Body.Bytes(), v); err != nil {
			t.Fatal(err)
		}
	}

	return recorder
}
