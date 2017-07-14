package system

import (
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

	s.Layer0.ScaleService(serviceID, 3)
	testutils.WaitFor(t, time.Second*10, time.Minute*5, func() bool {
		log.Debugf("Waiting for service to scale up")
		service := s.Layer0.GetService(serviceID)
		return service.RunningCount == 3
	})

	s.Layer0.ScaleService(serviceID, 1)
	testutils.WaitFor(t, time.Second*10, time.Minute*5, func() bool {
		log.Debugf("Waiting for service to scale down")
		service := s.Layer0.GetService(serviceID)
		return service.RunningCount == 1
	})
}
