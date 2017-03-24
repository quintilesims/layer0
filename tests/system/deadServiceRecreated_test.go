package system

import (
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/tests/system/clients"
	"testing"
	"time"
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

	testutils.WaitFor(t, time.Second*10, time.Minute, func() bool {
		logrus.Printf("Waiting for service to die")
		service := s.Layer0.GetService(serviceID)
		return service.RunningCount == 0
	})

	testutils.WaitFor(t, time.Second*10, time.Minute*2, func() bool {
		logrus.Printf("Waiting for service to recreate")
		service := s.Layer0.GetService(serviceID)
		return service.RunningCount == 1
	})
}
