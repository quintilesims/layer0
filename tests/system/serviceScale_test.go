package system

import (
	"log"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

// Test Resources:
// This test creates an environment named 'ss' that has a
// SystemTestService named 'sts'
func TestServiceScale(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/service_scale", nil)
	s.Terraform.Init()
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	serviceID := s.Terraform.Output("service_id")
	scale := 3

	req := models.UpdateServiceRequest{
		Scale: &scale,
	}

	s.Layer0.UpdateService(serviceID, req)
	testutils.WaitFor(t, time.Second*10, time.Minute*5, func() bool {
		log.Printf("[DEBUG] Waiting for service to scale up")
		service := s.Layer0.ReadService(serviceID)
		return service.RunningCount == 3
	})

	scale = 1

	req = models.UpdateServiceRequest{
		Scale: &scale,
	}

	s.Layer0.UpdateService(serviceID, req)
	testutils.WaitFor(t, time.Second*10, time.Minute*5, func() bool {
		log.Printf("[DEBUG] Waiting for service to scale down")
		service := s.Layer0.ReadService(serviceID)
		return service.RunningCount == 1
	})
}
