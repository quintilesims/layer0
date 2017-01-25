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

	// wait up to 3 minutes for the service to scale down
	waitFor(t, "Service to scale down", time.Minute*3, func() bool {
		svc := c.GetService("rssd", "sts")
		return svc.RunningCount == 1
	})

	if err := c.Client.RunRightSizer(); err != nil {
		t.Fatal(err)
	}

	// wait up to 5 minutes for the cluster to scale down
	waitFor(t, "Cluster to Scale Down", time.Minute*5, func() bool {
		env := c.GetEnvironment("rssd")
		return env.ClusterCount == 1
	})
}
