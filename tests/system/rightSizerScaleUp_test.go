package system

import (
	"testing"
)

// todo: use terraform modules to reduce layer0.tf duplication
func TestRightSizerScaleUp(t *testing.T) {
	c := startSystemTest(t, "cases/right_sizer_scale_up", nil)
	defer c.Destroy()

	/*
		env := c.GetEnvironment("rssu")
		svc := c.GetService(env.EnvironmentID, "baxter")
		//lb := c.GetLoadBalancer(env.EnvironmentID, "baxter")

		t.Log("Waiting up to 5 minutes for service to be running")
		waitFor(t, time.Minute*5, func() bool {
			print("waiting for service to be running")
			svc = c.GetService("", svc.ServiceID)
			return svc.RunningCount == 1
		})
	*/
}
