package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/retry"
	"github.com/zpatrick/rclient"
)

const MaxRetries = 3

type Config struct {
	Endpoint  string
	Token     string
	VerifySSL bool
}

type APIClient struct {
	client *rclient.RestClient
}

func NewAPIClient(c Config) *APIClient {
	httpClient := http.DefaultClient
	if !c.VerifySSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		httpClient = &http.Client{Transport: tr}
	}

	debugDoer := wrapDebugRequestDoer(httpClient)
	doer := wrapRetryRequestDoer(debugDoer)
	reader := newResponseReader()
	addAuthHeader := rclient.Header("Authorization", fmt.Sprintf("Basic %s", c.Token))

	restClient := rclient.NewRestClient(
		strings.TrimSuffix(c.Endpoint, "/"),
		rclient.Doer(doer),
		rclient.Reader(reader),
		rclient.RequestOptions(addAuthHeader))

	return &APIClient{
		client: restClient,
	}
}

func wrapRetryRequestDoer(doer rclient.RequestDoer) rclient.RequestDoerFunc {
	return func(req *http.Request) (*http.Response, error) {
		var response *http.Response
		var err error
		fn := func() (shouldRetry bool) {
			var resp *http.Response
			resp, err = doer.Do(req)
			if err != nil {
				return false
			}

			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				err = readServerError(resp)
				if err, ok := err.(*errors.ServerError); ok && err.Code == errors.EventualConsistencyError {
					log.Printf("[DEBUG] Client encountered eventual consistency error, will retry: %v", err)
					return true
				}

				return false
			}

			response = resp
			return false
		}

		if err := retry.Retry(fn, retry.WithMaxAttempts(MaxRetries)); err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		return response, nil
	}
}

func wrapDebugRequestDoer(doer rclient.RequestDoer) rclient.RequestDoer {
	return rclient.RequestDoerFunc(func(req *http.Request) (*http.Response, error) {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Request:\n%q", requestDump)

		resp, err := doer.Do(req)
		if err != nil {
			return nil, err
		}

		responseDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Response:\n%q", responseDump)

		return resp, nil
	})
}

func newResponseReader() rclient.ResponseReader {
	return func(resp *http.Response, v interface{}) error {
		defer resp.Body.Close()

		switch {
		case resp.StatusCode == 401:
			return fmt.Errorf("Invalid Auth Token. Have you tried running `l0-setup endpoint <instance>`?")
		case resp.StatusCode < 200, resp.StatusCode > 299:
			return readServerError(resp)
		case v == nil:
			return nil
		default:
			if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
				return fmt.Errorf("Failed to marshal response from Layer0 API: %v", err)
			}
		}

		return nil
	}
}

func readServerError(resp *http.Response) *errors.ServerError {
	defer resp.Body.Close()

	var se models.ServerError
	if err := json.NewDecoder(resp.Body).Decode(&se); err != nil {
		log.Printf("[DEBUG] Failed to decode server error: %v", err)
		return errors.Newf(errors.UnexpectedError, "Layer0 API returned a non-200 status code: %d", resp.StatusCode)
	}

	return errors.FromModel(se)
}
