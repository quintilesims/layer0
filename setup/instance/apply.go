package instance

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (i *Instance) Apply() error {
	if err := i.assertExists(); err != nil {
		return err
	}

	if err := i.validateInputs(); err != nil {
		return err
	}

	if err := i.Terraform.Apply(i.Dir); err != nil {
		return err
	}

	endpoint, err := i.Output(OUTPUT_ENDPOINT)
	if err != nil {
		return err
	}

	return i.waitForHealthyAPI(endpoint, time.Minute*10)
}

func (i *Instance) validateInputs() error {
	return nil
}

func (i *Instance) waitForHealthyAPI(endpoint string, timeout time.Duration) error {
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second * 15) {
		log.Printf("Waiting for API Service to be healthy... (%v)\n", time.Since(start))
		resp, err := http.Get(endpoint)
		if err != nil {
			log.Println("Error getting api: ", err)
			continue
		}

		defer resp.Body.Close()
		if code := resp.StatusCode; code < 200 || code > 299 {
			log.Println("API returned non-200 status code: %d", code)
			continue
		}

		return nil
	}

	return fmt.Errorf("API Service was not healthy after %v", timeout)
}
