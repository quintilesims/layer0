package system

import (
	"testing"
	"time"
)

func TestRightSizerScaleUp(t *testing.T) {
	c := startSystemTest(t, "cases/right_sizer_scale_up", nil)
	defer c.Destroy()

	svc := c.GetService("rssu", "sts")
	if _, err := c.Client.ScaleService(svc.ServiceID, 3); err != nil {
		t.Fatal(err)
	}

	// wait up to 5 minutes for the cluster to scale up
	waitFor(t, "Cluster to Scale Up", time.Minute, func() bool {
		env := c.GetEnvironment("rssu")
		return env.ClusterCount == 3
	})
}
