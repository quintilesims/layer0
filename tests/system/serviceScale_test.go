package system

import (
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
	"time"
)

// Test Resources:
// This test creates an environment named 'ss' that has a
// SystemTestService named 'sts'
func TestServiceScale(t *testing.T) {
	t.Parallel()

	test := NewSystemTest(t, "cases/service_scale", nil)
	test.Terraform.Apply()
	defer test.Terraform.Destroy()

	serviceID := test.Terraform.Output("service_id")

	test.L0Client.ScaleService(serviceID, 3)
	testutils.WaitFor(t, "Service to scale up", time.Minute*5, func() bool {
		logrus.Printf("Waiting for service to scale up")
		service := test.L0Client.GetService(serviceID)
		return service.RunningCount == 3
	})

	test.L0Client.ScaleService(serviceID, 1)
	testutils.WaitFor(t, "Service to scale down", time.Minute*5, func() bool {
		logrus.Printf("Waiting for service to scale down")
		service := test.L0Client.GetService(serviceID)
		return service.RunningCount == 1
	})
}
