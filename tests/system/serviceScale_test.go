package system

import (
	"log"
	"testing"
	"time"
)

// Test Resources:
// This test creates an environment named 'ss' that has a
// SystemTestService named 'sts'
func TestServiceScale(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/service_scale", nil)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	serviceID := s.Terraform.Output("service_id")
	deployID := s.Terraform.Output("deploy_id")

	s.Layer0.UpdateService(serviceID, deployID, 3)
	log.Printf("Waiting for service to scale up")
	service := s.Layer0.ReadService(serviceID)
	for start := time.Now(); time.Since(start) < time.Minute*5; time.Sleep(time.Second * 10) {
		if service.RunningCount == 3 {
			continue
		}
	}

	if service.RunningCount != 3 {
		t.Fatalf("[ERROR] Timeout reached after %v", time.Minute*5)
	}

	s.Layer0.UpdateService(serviceID, deployID, 1)
	log.Printf("Waiting for service to scale down")
	for start := time.Now(); time.Since(start) < time.Minute*5; time.Sleep(time.Second * 10) {
		if service.RunningCount == 1 {
			continue
		}
	}

	if service.RunningCount != 1 {
		t.Fatalf("[ERROR] Timeout reached after %v", time.Minute*5)
	}
}
