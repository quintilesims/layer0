package rclient

import (
	"net/http"
)

// A RequestDoer sends a *http.Request and returns a *http.Response.
type RequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// A RequestDoerFunc is a function that implements the RequestDoer interface.
type RequestDoerFunc func(*http.Request) (*http.Response, error)

// Do executes the RequestDoerFunc.
func (d RequestDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return d(req)
}
