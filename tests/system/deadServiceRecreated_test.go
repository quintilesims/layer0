package system

import (
	"testing"
	"time"
)

func TestDeadServiceRecreated(t *testing.T) {
	c := startSystemTest(t, "cases/dead_service_recreated", nil)
	defer c.Destroy()

	s := c.GetSystemTestService("dsr", "sts")
	s.Die()

	// wait up to 1 minute for the service to die
	waitFor(t, "Service to Die", time.Minute, func() bool {
		svc := c.GetService("dsr", "sts")
		return svc.RunningCount == 0
	})

	// wait up to 2 minutes for the service to be recreated
	waitFor(t, "Service to Recreate", time.Minute*2, func() bool {
		svc := c.GetService("dsr", "sts")
		return svc.RunningCount == 1
	})
}
