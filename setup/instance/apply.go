package instance

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/quintilesims/layer0/common/config"
)

func (l *LocalInstance) Apply(wait bool) error {
	if err := l.assertExists(); err != nil {
		return err
	}

	if err := l.Terraform.Apply(l.Dir); err != nil {
		return err
	}

	if wait {
		endpoint, err := l.Output(config.FlagEndpoint.GetName())
		if err != nil {
			return err
		}

		token, err := l.Output(config.FlagToken.GetName())
		if err != nil {
			return err
		}

		return l.waitForHealthyAPI(endpoint, token, time.Minute*10)
	}

	return nil
}

func (l *LocalInstance) waitForHealthyAPI(endpoint, token string, timeout time.Duration) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+token)
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second * 15) {
		log.Printf("[INFO] Waiting for API Service to be healthy... (%s)", time.Since(start).String())

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[WARN] Error occurred during GET %s: %v", endpoint, err)
			continue
		}

		defer resp.Body.Close()
		if code := resp.StatusCode; code < 200 || code > 299 {
			log.Printf("[WARN] API returned non-200 status code: %d", code)
			continue
		}

		return nil
	}

	return fmt.Errorf("API Service was not healthy after %v", timeout)
}
