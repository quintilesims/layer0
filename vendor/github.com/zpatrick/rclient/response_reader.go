package rclient

import (
	"encoding/json"
	"net/http"
)

// A ResponseReader attempts to read a *http.Response into v.
type ResponseReader func(resp *http.Response, v interface{}) error

// ReadJSONResponse attempts to marshal the response body into v
// if and only if the response StatusCode is in the 200 range.
// Otherwise, an error is thrown.
// It assumes the response body is in JSON format.
func ReadJSONResponse(resp *http.Response, v interface{}) error {
	switch {
	case resp.StatusCode < 200, resp.StatusCode > 299:
		return NewResponseErrorf(resp, "Invalid status code: %d", resp.StatusCode)
	case v == nil:
		return nil
	default:
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return NewResponseError(resp, err.Error())
		}

		resp.Body.Close()
	}

	return nil
}
