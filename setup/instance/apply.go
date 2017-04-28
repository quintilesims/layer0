package instance

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"net/http"
	"time"
)

func (l *LocalInstance) Apply(wait bool) error {
	if err := l.assertExists(); err != nil {
		return err
	}

	if err := l.Terraform.Apply(l.Dir); err != nil {
		return err
	}

	endpoint, err := l.Output(OUTPUT_ENDPOINT)
	if err != nil {
		return err
	}

	if wait {
		return l.waitForHealthyAPI(endpoint, time.Minute*10)
	}

	return nil
}

func (l *LocalInstance) waitForHealthyAPI(endpoint string, timeout time.Duration) error {
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second * 15) {
		logrus.Infof("Waiting for API Service to be healthy... (%3s)", time.Since(start))

		resp, err := http.Get(endpoint)
		if err != nil {
			logrus.Debugf("Error occurred during GET %s: %v", endpoint, err)
			continue
		}

		defer resp.Body.Close()
		if code := resp.StatusCode; code < 200 || code > 299 {
			logrus.Debugf("API returned non-200 status code: %d", code)
			continue
		}

		return nil
	}

	return fmt.Errorf("API Service was not healthy after %v", timeout)
}
