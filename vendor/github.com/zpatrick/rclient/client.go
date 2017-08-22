package rclient

import (
	"fmt"
	"net/http"
)

// RestClient builds, executes, and reads http requests/responses.
type RestClient struct {
	Host           string
	RequestBuilder RequestBuilder
	RequestDoer    RequestDoer
	ResponseReader ResponseReader
	RequestOptions []RequestOption
}

// NewRestClient returns a new RestClient with all of the default fields.
// Any of the default fields can be changed with the options param.
func NewRestClient(host string, options ...ClientOption) (*RestClient, error) {
	r := &RestClient{
		Host:           host,
		RequestBuilder: BuildJSONRequest,
		RequestDoer:    http.DefaultClient,
		ResponseReader: ReadJSONResponse,
		RequestOptions: []RequestOption{},
	}

	for _, option := range options {
		if err := option(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

// Delete passes its params to RestClient.Do() with the "DELETE" method.
func (r *RestClient) Delete(path string, body, v interface{}, options ...RequestOption) error {
	return r.Do("DELETE", path, body, v, options...)
}

// Get passes its params to RestClient.Do() with the "GET" method.
func (r *RestClient) Get(path string, v interface{}, options ...RequestOption) error {
	return r.Do("GET", path, nil, v, options...)
}

// Patch passes its params to RestClient.Do() with the "PATCH" method.
func (r *RestClient) Patch(path string, body, v interface{}, options ...RequestOption) error {
	return r.Do("PATCH", path, body, v, options...)
}

// Post passes its params to RestClient.Do() with the "POST" method.
func (r *RestClient) Post(path string, body, v interface{}, options ...RequestOption) error {
	return r.Do("POST", path, body, v, options...)
}

// Put passes its params to RestClient.Do() with the "PUT" method.
func (r *RestClient) Put(path string, body, v interface{}, options ...RequestOption) error {
	return r.Do("PUT", path, body, v, options...)
}

// Do orchestrates building, performing, and reading http requests and responses.
func (r *RestClient) Do(method, path string, body, v interface{}, options ...RequestOption) error {
	url := fmt.Sprintf("%s%s", r.Host, path)
	options = append(r.RequestOptions, options...)

	req, err := r.RequestBuilder(method, url, body, options...)
	if err != nil {
		return err
	}

	resp, err := r.RequestDoer.Do(req)
	if err != nil {
		return err
	}

	return r.ResponseReader(resp, v)
}
