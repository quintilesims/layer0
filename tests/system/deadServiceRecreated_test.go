package system

import (
	"log"
	"testing"
	"time"

	"github.com/quintilesims/layer0/tests/clients"
)

// Test Resources:
// This test creates an environment named 'dsr' that has a
// SystemTestService named 'sts'
func TestDeadServiceRecreated(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/dead_service_recreated", nil)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	serviceID := s.Terraform.Output("service_id")
	serviceURL := s.Terraform.Output("service_url")

	sts := clients.NewSTSTestClient(t, serviceURL)
	sts.WaitForHealthy(time.Minute * 3)
	sts.SetHealth("die")

	log.Printf("[DEBUG] Waiting for service to die")
	service := s.Layer0.ReadService(serviceID)
	for start := time.Now(); time.Since(start) < time.Minute*2; time.Sleep(time.Second * 10) {
		if service.RunningCount == 0 {
			continue
		}
	}

	if service.RunningCount != 0 {
		t.Fatalf("[ERROR] Timeout reached after %v", time.Minute*2)
	}

	log.Printf("[DEBUG] Waiting for service to recreate")
	for start := time.Now(); time.Since(start) < time.Minute*2; time.Sleep(time.Second * 10) {
		if service.RunningCount == 1 {
			continue
		}
	}

	if service.RunningCount != 1 {
		t.Fatalf("[ERROR] Timeout reached after %v", time.Minute*2)
	}
}
