package system

import (
	"github.com/quintilesims/layer0/common/testutils"
	"net/http"
	"testing"
	"time"
)

// Test Resources:
// This test creates an environment named 'ws' that has a
// Windows service named 'windows' running in iis
func TestWindowsService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestWindowsService in short mode")
	}

	t.Parallel()

	s := NewSystemTest(t, "cases/windows_service", nil)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	serviceID := s.Terraform.Output("service_id")
	serviceURL := s.Terraform.Output("service_url")

	testutils.WaitFor(t, time.Second*30, time.Minute*45, func() bool {
		log.Debugf("Waiting for windows service to run")
		service := s.Layer0.GetService(serviceID)
		return service.RunningCount == 1
	})

	testutils.WaitFor(t, time.Second*30, time.Minute*10, func() bool {
		log.Debugf("Waiting for service to be healthy")
		resp, err := http.Get(serviceURL)
		if err != nil {
			log.Debugf("There was an error checking the window service's url: %v", err)
			return false
		}

		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Debugf("Windows service returned non-200 status: %d", resp.StatusCode)
			return false
		}

		return true
	})
}
