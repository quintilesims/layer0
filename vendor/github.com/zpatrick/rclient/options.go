package rclient

import (
	"net/http"
	"net/url"
)

// A ClientOption configures a *RestClient.
type ClientOption func(client *RestClient) error

// Builder sets the RequestBuilder field of a RestClient.
func Builder(builder RequestBuilder) ClientOption {
	return func(r *RestClient) error {
		r.RequestBuilder = builder
		return nil
	}
}

// Doer sets the RequestDoer field of a RestClient.
func Doer(doer RequestDoer) ClientOption {
	return func(r *RestClient) error {
		r.RequestDoer = doer
		return nil
	}
}

// Reader sets the ResponseReader field of a RestClient.
func Reader(reader ResponseReader) ClientOption {
	return func(r *RestClient) error {
		r.ResponseReader = reader
		return nil
	}
}

// RequestOptions sets the RequestOptions field of a RestClient.
func RequestOptions(options ...RequestOption) ClientOption {
	return func(r *RestClient) error {
		r.RequestOptions = append(r.RequestOptions, options...)
		return nil
	}
}

// A RequestOption configures a *http.Request.
type RequestOption func(req *http.Request) error

// BasicAuth adds the specified username and password as basic auth to a request.
func BasicAuth(user, pass string) RequestOption {
	return func(req *http.Request) error {
		req.SetBasicAuth(user, pass)
		return nil
	}
}

// Header adds the specified name and value as a header to a request.
func Header(name, val string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Add(name, val)
		return nil
	}
}

// Headers adds the specified names and values as headers to a request
func Headers(headers map[string]string) RequestOption {
	return func(req *http.Request) error {
		for name, val := range headers {
			req.Header.Add(name, val)
		}

		return nil
	}
}

// Query adds the specified query to a request.
func Query(query url.Values) RequestOption {
	return func(req *http.Request) error {
		req.URL.RawQuery = query.Encode()
		return nil
	}
}
