package system

import (
	"log"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/testutils"
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
	testutils.WaitFor(t, time.Second*10, time.Minute*5, func() bool {
		log.Printf("[DEBUG] Waiting for service to scale up")
		service := s.Layer0.ReadService(serviceID)
		return service.RunningCount == 3
	})

	s.Layer0.UpdateService(serviceID, deployID, 1)
	testutils.WaitFor(t, time.Second*10, time.Minute*5, func() bool {
		log.Printf("[DEBUG] Waiting for service to scale down")
		service := s.Layer0.ReadService(serviceID)
		return service.RunningCount == 1
	})
}
