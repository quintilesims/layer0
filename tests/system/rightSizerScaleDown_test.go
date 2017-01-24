package system

import (
	"testing"
	"time"
)

func TestRightSizerScaleDown(t *testing.T) {
	c := startSystemTest(t, "cases/right_sizer_scale_down", nil)
	defer c.Destroy()

	svc := c.GetService("rssd", "sts")
	if _, err := c.Client.ScaleService(svc.ServiceID, 1); err != nil {
		t.Fatal(err)
	}

	// todo: trigger api right sizer
	// c.Client.RunRightSizer()

	// wait up to 5 minutes for the cluster to scale down
	waitFor(t, "Cluster to Scale Down", time.Minute, func() bool {
		env := c.GetEnvironment("rssd")
		return env.ClusterCount == 1
	})
}
