package system

import (
	"log"
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

	service := s.Layer0.ReadService(serviceID)
	log.Printf("Waiting for windows service to run")
	for start := time.Now(); time.Since(start) < time.Minute*45; time.Sleep(time.Second * 30) {
		if service.RunningCount == 1 {
			continue
		}
	}

	if service.RunningCount != 1 {
		t.Fatalf("[ERROR] Timeout reached after %v", time.Minute*45)
	}

	log.Printf("Waiting for service to be healthy")
	for start := time.Now(); time.Since(start) < time.Minute*10; time.Sleep(time.Second * 30) {
		resp, err := http.Get(serviceURL)
		if err != nil {
			t.Fatalf("[ERROR] There was an error checking the Windows service's URL: %v", err)
		}

		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			t.Fatalf("[ERROR] Windows service returned non-200 status: %d", resp.StatusCode)
		}
	}
}
