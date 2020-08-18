package rclient

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// A RequestBuilder creates a *http.Request from the given parameters.
// It is important that each option gets added to the generated request:
//  req, _ := http.NewRequest(...)
//  for _, option := range options
//      if err := option(req); err != nil {
//          return nil, err
//      }
//  }
type RequestBuilder func(method, url string, body interface{}, options ...RequestOption) (*http.Request, error)

// BuildJSONRequest creates a new *http.Request with the specified method, url and body in JSON format.
func BuildJSONRequest(method, url string, body interface{}, options ...RequestOption) (*http.Request, error) {
	b := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(b).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}

	if b.Len() > 0 {
		req.Header.Add("content-type", "application/json")
	}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
