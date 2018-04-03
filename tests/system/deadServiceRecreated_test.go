package system

import (
	"log"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/testutils"
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

	statelessServiceID := s.Terraform.Output("stateless_service_id")
	statelessServiceURL := s.Terraform.Output("stateless_service_url")

	statelessClient := clients.NewSTSTestClient(t, statelessServiceURL)
	statelessClient.WaitForHealthy(time.Minute * 3)
	statelessClient.SetHealth("die")

	testutils.WaitFor(t, time.Second*10, time.Minute, func() bool {
		log.Printf("[DEBUG] Waiting for stateless service to die")
		service := s.Layer0.ReadService(statelessServiceID)
		return service.RunningCount == 0
	})

	testutils.WaitFor(t, time.Second*10, time.Minute*2, func() bool {
		log.Printf("[DEBUG] Waiting for stateless service to recreate")
		service := s.Layer0.ReadService(statelessServiceID)
		return service.RunningCount == 1
	})

	statefulServiceID := s.Terraform.Output("stateful_service_id")
	statefulServiceURL := s.Terraform.Output("stateful_service_url")

	statefulClient := clients.NewSTSTestClient(t, statefulServiceURL)
	statefulClient.WaitForHealthy(time.Minute * 3)
	statefulClient.SetHealth("die")

	testutils.WaitFor(t, time.Second*10, time.Minute, func() bool {
		log.Printf("[DEBUG] Waiting for stateful service to die")
		service := s.Layer0.ReadService(statefulServiceID)
		return service.RunningCount == 0
	})

	testutils.WaitFor(t, time.Second*10, time.Minute*2, func() bool {
		log.Printf("[DEBUG] Waiting for stateful service to recreate")
		service := s.Layer0.ReadService(statefulServiceID)
		return service.RunningCount == 1
	})
}
