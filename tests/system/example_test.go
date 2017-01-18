package system

import (
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestOne(t *testing.T) {
	c := runSystemTest(t, "example")
	defer c.Destroy()

	demo := c.GetEnvironment("demo")
	testutils.AssertEqual(t, demo.EnvironmentName, "some bogus name")
}

func TestTwo(t *testing.T) {
	c := runSystemTest(t, "example")
	defer c.Destroy()

	apiENV := c.GetEnvironment("api")
	testutils.AssertEqual(t, apiENV.EnvironmentName, "api")

	apiLB := c.GetLoadBalancer("api")
	testutils.AssertEqual(t, apiLB.LoadBalancerName, "api")

	apiSVC := c.GetService("api")
	testutils.AssertEqual(t, apiSVC.ServiceName, "api-svc")
}
