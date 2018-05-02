package client

import (
	"io"
	"io/ioutil"
	"time"

	log "github.com/Sirupsen/logrus"

	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/dghubble/sling"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/retry"
	"github.com/quintilesims/layer0/common/waitutils"
)

type DoerFunc func(req *http.Request) (*http.Response, error)

func (d DoerFunc) Do(req *http.Request) (*http.Response, error) {
	return d(req)
}

func logSling(httpClient *http.Client) sling.Doer {
	readBody := func(body io.Reader) (io.ReadCloser, string) {
		bodyBytes, _ := ioutil.ReadAll(body)
		original := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return original, string(bodyBytes)
	}

	return DoerFunc(func(req *http.Request) (*http.Response, error) {
		var requestBody string
		if req.Body != nil {
			req.Body, requestBody = readBody(req.Body)
		}

		log.Debugf("Sent: %s %s %s\n", req.Method, req.URL, requestBody)
		resp, err := httpClient.Do(req)

		var responseBody string
		if resp != nil {
			if resp.Body != nil {
				resp.Body, responseBody = readBody(resp.Body)
			}

			log.Debugf("Received: %s %s \n", resp.Status, responseBody)
		}

		return resp, err
	})
}

func retrySling(httpClient *http.Client) sling.Doer {
	readBody := func(body io.Reader) (io.ReadCloser, string) {
		bodyBytes, _ := ioutil.ReadAll(body)
		original := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return original, string(bodyBytes)
	}

	return DoerFunc(func(req *http.Request) (*http.Response, error) {
		if req.Body != nil {
			req.Body, _ = readBody(req.Body)
		}

		var resp *http.Response
		retryFN := func() (shouldRetry bool, err error) {
			if resp, err = httpClient.Do(req); err != nil {
				if strings.Contains(err.Error(), "EOF") {
					return true, nil
				}

				return false, err
			}

			return false, nil
		}

		err := retry.Retry(retryFN, retry.WithTimeout(time.Minute*3), retry.WithDelay(time.Second*5))

		return resp, err
	})
}

type Config struct {
	Endpoint      string
	Token         string
	VerifySSL     bool
	VerifyVersion bool
	Clock         waitutils.Clock
}

type APIClient struct {
	Endpoint      string
	Token         string
	VerifyVersion bool
	Clock         waitutils.Clock
	httpClient    *http.Client
	once          sync.Once
}

func NewAPIClient(config Config) *APIClient {
	var httpClient *http.Client
	if !config.VerifySSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		httpClient = &http.Client{Transport: tr}
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if config.Clock == nil {
		config.Clock = waitutils.RealClock{}
	}

	return &APIClient{
		httpClient:    httpClient,
		Endpoint:      config.Endpoint,
		Token:         config.Token,
		VerifyVersion: config.VerifyVersion,
		Clock:         config.Clock,
	}
}

func (c *APIClient) Sling(path string) *sling.Sling {
	return sling.New().
		Client(c.httpClient).
		Base(c.Endpoint).
		Path(path).
		Set("Authorization", fmt.Sprintf("Basic %s", c.Token)).
		Doer(logSling(c.httpClient)).
		Doer(retrySling(c.httpClient))
}

func (c *APIClient) Execute(sling *sling.Sling, receive interface{}) error {
	if _, err := c.execute(sling, receive); err != nil {
		return err
	}

	return nil
}

func (c *APIClient) ExecuteWithJob(sling *sling.Sling) (string, error) {
	var response *string
	resp, err := c.execute(sling, &response)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusAccepted {
		if jobID := resp.Header.Get("X-JobID"); jobID != "" {
			return jobID, nil
		}

		return "", fmt.Errorf("API returned %v but no 'X-JobID' header", http.StatusAccepted)
	}

	return "", fmt.Errorf("Failed to get job from response: Status was %v (expected %v)", resp.StatusCode, http.StatusAccepted)
}

func (c *APIClient) execute(sling *sling.Sling, receive interface{}) (*http.Response, error) {
	var serverError *ServerError
	resp, err := sling.Receive(receive, &serverError)

	if err != nil {
		if strings.Contains(err.Error(), "x509: certificate is valid for") {
			return nil, sslError(err)
		}

		if resp != nil && resp.StatusCode == 401 {
			return nil, fmt.Errorf("Invalid Auth Token. Have you tried running `l0-setup endpoint <prefix>`?")
		}

		if _, ok := err.(*url.Error); ok {
			return nil, fmt.Errorf("Unable to connect to API with error: %v", err)
		}

		return nil, err
	}

	if serverError != nil {
		return nil, serverError.ToCommonError()
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Layer0 API returned invalid status code: %s", resp.Status)
	}

	if err := c.verifyVersion(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *APIClient) verifyVersion(resp *http.Response) error {
	if !c.VerifyVersion {
		return nil
	}

	if cli, api := config.CLIVersion(), resp.Header.Get("Version"); cli != api {
		text := fmt.Sprintf("API and CLI version mismatch (CLI: '%s', API: '%s')\n", cli, api)
		text += fmt.Sprintf("To disable this warning, set %s=\"1\"", config.SKIP_VERSION_VERIFY)
		return fmt.Errorf(text)
	}

	return nil
}
