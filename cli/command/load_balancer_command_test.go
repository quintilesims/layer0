package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/urfave/cli"
)

func TestParsePort(t *testing.T) {
	cases := []struct {
		Target      string
		Certificate string
		Expected    models.Port
	}{
		{
			Target: "80:80/tcp",
			Expected: models.Port{
				HostPort:      80,
				ContainerPort: 80,
				Protocol:      "tcp",
			},
		},
		{
			Target: "80:80/http",
			Expected: models.Port{
				HostPort:      80,
				ContainerPort: 80,
				Protocol:      "http",
			},
		},
		{
			Target: "8080:80/http",
			Expected: models.Port{
				HostPort:      8080,
				ContainerPort: 80,
				Protocol:      "http",
			},
		},
		{
			Target:      "80:80/https",
			Certificate: "crt_name",
			Expected: models.Port{
				HostPort:        80,
				ContainerPort:   80,
				Protocol:        "https",
				CertificateName: "crt_name",
			},
		},
		{
			Target:      "80:80/https",
			Certificate: "arn:aws:iam::12345:server-certificate/crt_name",
			Expected: models.Port{
				HostPort:       80,
				ContainerPort:  80,
				Protocol:       "https",
				CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name",
			},
		},
	}

	for _, c := range cases {
		result, err := parsePort(c.Target, c.Certificate)
		if err != nil {
			t.Fatal(err)
		}

		testutils.AssertEqual(t, *result, c.Expected)
	}
}

func TestParsePortErrors(t *testing.T) {
	cases := map[string]string{
		"Missing HOST_PORT":          ":80/tcp",
		"Missing CONTAINER_PORT":     "80:/tcp",
		"Missing PROTOCOL":           "80:80",
		"Non-integer HOST_PORT":      "80p:80/tcp",
		"Non-integer CONTAINER_PORT": "80:80p/tcp",
	}

	for name, input := range cases {
		if _, err := parsePort(input, ""); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestLoadBalancerAddPort(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("load_balancer", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetLoadBalancer("id").
		Return(&models.LoadBalancer{}, nil)

	port := models.Port{
		HostPort:        443,
		ContainerPort:   80,
		Protocol:        "https",
		CertificateName: "cert_name",
	}

	tc.Client.EXPECT().
		UpdateLoadBalancerPorts("id", []models.Port{port}).
		Return(&models.LoadBalancer{}, nil)

	flags := map[string]interface{}{"certificate": "cert_name"}
	c := testutils.GetCLIContext(t, []string{"name", "443:80/https"}, flags)
	if err := command.AddPort(c); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerAddPort_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
		"Missing PORT arg": testutils.GetCLIContext(t, []string{"name"}, nil),
	}

	for name, c := range contexts {
		if err := command.AddPort(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestCreateLoadBalancer(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "environment").
		Return([]string{"environmentID"}, nil)

	healthCheck := models.HealthCheck{
		Target:             "TCP:80",
		Interval:           30,
		Timeout:            5,
		HealthyThreshold:   10,
		UnhealthyThreshold: 2,
	}

	ports := []models.Port{
		{
			HostPort:        443,
			ContainerPort:   80,
			Protocol:        "https",
			CertificateName: "cert_name",
		},
		{
			HostPort:        8000,
			ContainerPort:   8000,
			Protocol:        "http",
			CertificateName: "",
		},
	}

	tc.Client.EXPECT().
		CreateLoadBalancer("name", "environmentID", healthCheck, ports, false, 60).
		Return(&models.LoadBalancer{}, nil)

	flags := map[string]interface{}{
		"port":                            []string{"443:80/https", "8000:8000/http"},
		"certificate":                     "cert_name",
		"private":                         true,
		"healthcheck-target":              "TCP:80",
		"healthcheck-interval":            30,
		"healthcheck-timeout":             5,
		"healthcheck-healthy-threshold":   10,
		"healthcheck-unhealthy-threshold": 2,
		"idle-timeout":                    60,
	}

	c := testutils.GetCLIContext(t, []string{"environment", "name"}, flags)
	if err := command.Create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateLoadBalancer_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": testutils.GetCLIContext(t, nil, nil),
		"Missing NAME arg":        testutils.GetCLIContext(t, []string{"environment"}, nil),
	}

	for name, c := range contexts {
		if err := command.Create(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestDeleteLoadBalancer(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("load_balancer", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		DeleteLoadBalancer("id").
		Return("jobid", nil)

	c := testutils.GetCLIContext(t, []string{"name"}, nil)
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteLoadBalancerWait(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("load_balancer", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		DeleteLoadBalancer("id").
		Return("jobid", nil)

	tc.Client.EXPECT().
		WaitForJob("jobid", testutils.TEST_TIMEOUT).
		Return(nil)

	c := testutils.GetCLIContext(t, []string{"name"}, map[string]interface{}{"wait": true})
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteLoadBalancer_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Delete(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestGetLoadBalancer(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("load_balancer", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetLoadBalancer("id").
		Return(&models.LoadBalancer{}, nil)

	c := testutils.GetCLIContext(t, []string{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetLoadBalancer_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Get(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestListLoadBalancers(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	tc.Client.EXPECT().
		ListLoadBalancers().
		Return([]*models.LoadBalancerSummary{}, nil)

	c := testutils.GetCLIContext(t, nil, nil)
	if err := command.List(c); err != nil {
		t.Fatal(err)
	}
}

func TestHealthCheck_noUpdateRequired(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("load_balancer", "env").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetLoadBalancer("id").
		Return(&models.LoadBalancer{}, nil)

	c := testutils.GetCLIContext(t, []string{"env", "name"}, nil)
	if err := command.HealthCheck(c); err != nil {
		t.Fatal(err)
	}
}

func TestHealthCheck_partialUpdateRequired(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	existingHealthCheck := models.HealthCheck{
		Target:             "TCP:80",
		Interval:           30,
		Timeout:            5,
		HealthyThreshold:   2,
		UnhealthyThreshold: 2,
	}

	tc.Resolver.EXPECT().
		Resolve("load_balancer", "env").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetLoadBalancer("id").
		Return(&models.LoadBalancer{
			HealthCheck: existingHealthCheck,
		}, nil)

	expectedHealthCheck := models.HealthCheck{
		Target:             "TCP:88",
		Interval:           45,
		Timeout:            5,
		HealthyThreshold:   2,
		UnhealthyThreshold: 2,
	}

	tc.Client.EXPECT().
		UpdateLoadBalancerHealthCheck("id", expectedHealthCheck)

	flags := map[string]interface{}{
		"set-target":   "TCP:88",
		"set-interval": 45,
	}

	c := testutils.GetCLIContext(t, []string{"env", "name"}, flags)
	if err := command.HealthCheck(c); err != nil {
		t.Fatal(err)
	}
}

func TestHealthCheck_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Non-int '--set-interval' flag":            testutils.GetCLIContext(t, []string{"name"}, map[string]interface{}{"set-interval": "two"}),
		"Non-int '--set-timeout' flag":             testutils.GetCLIContext(t, []string{"name"}, map[string]interface{}{"set-timeout": "two"}),
		"Non-int '--set-healthy-threshold' flag":   testutils.GetCLIContext(t, []string{"name"}, map[string]interface{}{"set-healthy-threshold": "two"}),
		"Non-int '--set-unhealthy-threshold' flag": testutils.GetCLIContext(t, []string{"name"}, map[string]interface{}{"set-unhealthy-threshold": "two"}),
		"Missing NAME arg":                         testutils.GetCLIContext(t, nil, map[string]interface{}{"set-interval": 2}),
	}

	for name, c := range contexts {
		if err := command.HealthCheck(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestLoadBalancerIdleTimeout(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("load_balancer", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetLoadBalancer("id").
		Return(&models.LoadBalancer{}, nil)

	idleTimeout := 60

	tc.Client.EXPECT().
		UpdateLoadBalancerIdleTimeout("id", idleTimeout).
		Return(&models.LoadBalancer{}, nil)

	c := testutils.GetCLIContext(t, nil, nil)
	if err := command.IdleTimeout(c); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerIdleTimeout_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg":    testutils.GetCLIContext(t, nil, nil),
		"Missing TIMEOUT arg": testutils.GetCLIContext(t, []string{"name"}, nil),
	}

	for name, c := range contexts {
		if err := command.IdleTimeout(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}
