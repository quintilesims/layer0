package system

import (
	//"github.com/quintilesims/layer0/common/testutils"
	"testing"
	"time"
)

// todo: use terraform modules to reduce layer0.tf duplication
func TestDeadServiceRecreated(t *testing.T) {
	c := startSystemTest(t, "cases/dead_service_recreated", nil)
	defer c.Destroy()

	/*
		env := c.GetEnvironment("dsr")
		svc := c.GetService(env.EnvironmentID, "baxter")
		lb := c.GetLoadBalancer(env.EnvironmentID, "baxter")

		t.Log("Waiting up to 5 minutes for service to be running")
		waitFor(t, time.Minute*5, func() bool {
			print("waiting for service to be running")
			svc = c.GetService("", svc.ServiceID)
			return svc.RunningCount == 1
		})

		t.Log("Telling service to die")
		b := NewBaxter(t, lb.URL)
		b.Die()

		t.Log("Waiting up to 1 minute for service to die")
		waitFor(t, time.Minute*1, func() bool {
			print("waiting for service to die")
			svc = c.GetService("", svc.ServiceID)
			return svc.RunningCount == 0
		})

		t.Log("Waiting up to 2 minute for the service to get recreated")
		waitFor(t, time.Minute*2, func() bool {
			print("waiting for service reacreate")
			svc = c.GetService("", svc.ServiceID)
			return svc.RunningCount == 1
		})

		print("test is done")
	*/
}

func waitFor(t *testing.T, timeout time.Duration, shouldStopWaiting func() bool) {
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second * 5) {
		if shouldStopWaiting() {
			return
		}
	}

	t.Fatalf("Timeout after %v", timeout)
}
