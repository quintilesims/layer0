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

	test := NewSystemTest(t, "cases/dead_service_recreated", nil)
	test.Terraform.Apply()
	defer test.Terraform.Destroy()

	serviceID := test.Terraform.Output("service_id")
	serviceURL := test.Terraform.Output("service_url")

	stsClient := clients.NewSTSTestClient(t, serviceURL)
	stsClient.WaitForHealthy(time.Minute * 3)
	stsClient.SetHealth("die")

	testutils.WaitFor(t, time.Minute, func() bool {
		logrus.Printf("Waiting for service to die")
		service := test.L0Client.GetService(serviceID)
		return service.RunningCount == 0
	})

	testutils.WaitFor(t, time.Minute*2, func() bool {
		logrus.Printf("Waiting for service to recreate")
		service := test.L0Client.GetService(serviceID)
		return service.RunningCount == 1
	})
}
