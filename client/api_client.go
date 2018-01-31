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
	"github.com/zpatrick/rclient"
)

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

	doer := newDebugRequestDoer(httpClient)
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

func newDebugRequestDoer(httpClient *http.Client) rclient.RequestDoer {
	return rclient.RequestDoerFunc(func(req *http.Request) (*http.Response, error) {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Request:\n%q", requestDump)

		resp, err := httpClient.Do(req)
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
			var se models.ServerError
			if err := json.NewDecoder(resp.Body).Decode(&se); err != nil {
				log.Printf("[DEBUG] Failed to decode server error: %v", err)
				return fmt.Errorf("Layer0 API returned a non-200 status code: %d", resp.StatusCode)
			}

			return errors.FromModel(se)
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
