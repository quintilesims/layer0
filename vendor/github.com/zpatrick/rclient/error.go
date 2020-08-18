package rclient

import (
	"fmt"
	"net/http"
)

type ResponseError struct {
	Response *http.Response
	Message  string
}

func (e *ResponseError) Error() string {
	return e.Message
}

func NewResponseError(resp *http.Response, message string) *ResponseError {
	return &ResponseError{
		Response: resp,
		Message:  message,
	}
}

func NewResponseErrorf(resp *http.Response, format string, tokens ...interface{}) *ResponseError {
	return NewResponseError(resp, fmt.Sprintf(format, tokens...))
}
