package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestLoadBalancerAddPort(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	base.Client.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(&models.LoadBalancer{}, nil)

	ports := []models.Port{
		{HostPort: 443, ContainerPort: 80, Protocol: "https", CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name"},
	}

	req := models.UpdateLoadBalancerRequest{
		Ports: &ports,
	}

	base.Client.EXPECT().
		UpdateLoadBalancer("lb_id", req).
		Return(nil)

	base.Client.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(&models.LoadBalancer{}, nil)

	flags := map[string]interface{}{
		"certificate": "arn:aws:iam::12345:server-certificate/crt_name",
	}

	c := testutils.NewTestContext(t, []string{"lb_name", "443:80/https"}, flags)
	if err := command.addPort(c); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerAddPortInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
		"Missing PORT arg": testutils.NewTestContext(t, []string{"lb_name"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.addPort(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestCreateLoadBalancer(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		IsPublic:         false,
		Ports: []models.Port{
			{HostPort: 443, ContainerPort: 80, Protocol: "https", CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name"},
			{HostPort: 22, ContainerPort: 22, Protocol: "tcp", CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name"},
		},
		HealthCheck: models.HealthCheck{
			Target:             "tcp:80",
			Interval:           5,
			Timeout:            6,
			HealthyThreshold:   7,
			UnhealthyThreshold: 8,
		},
	}

	base.Client.EXPECT().
		CreateLoadBalancer(req).
		Return("lb_id", nil)

	base.Client.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(&models.LoadBalancer{}, nil)

	flags := map[string]interface{}{
		"private": true,
		"port": []string{
			"443:80/https",
			"22:22/tcp",
		},
		"certificate":                     "arn:aws:iam::12345:server-certificate/crt_name",
		"healthcheck-target":              "tcp:80",
		"healthcheck-interval":            5,
		"healthcheck-timeout":             6,
		"healthcheck-healthy-threshold":   7,
		"healthcheck-unhealthy-threshold": 8,
	}

	c := testutils.NewTestContext(t, []string{"env_name", "lb_name"}, flags)
	if err := command.create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateLoadBalancerInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": testutils.NewTestContext(t, nil, nil),
		"Missing NAME arg":        testutils.NewTestContext(t, []string{"env_name"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.create(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestDeleteLoadBalancer(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	base.Client.EXPECT().
		DeleteLoadBalancer("lb_id").
		Return(nil)

	c := testutils.NewTestContext(t, []string{"lb_name"}, nil)
	if err := command.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteLoadBalancerInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.delete(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestLoadBalancerDropPort(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	loadBalancer := &models.LoadBalancer{
		Ports: []models.Port{
			{HostPort: 443},
		},
	}

	base.Client.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(loadBalancer, nil)

	ports := []models.Port{}
	req := models.UpdateLoadBalancerRequest{
		Ports: &ports,
	}

	base.Client.EXPECT().
		UpdateLoadBalancer("lb_id", req).
		Return(nil)

	base.Client.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(&models.LoadBalancer{}, nil)

	c := testutils.NewTestContext(t, []string{"lb_name", "443"}, nil)
	if err := command.dropPort(c); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerDropPortInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
		"Missing PORT arg": testutils.NewTestContext(t, []string{"lb_name"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.dropPort(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestListLoadBalancers(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	base.Client.EXPECT().
		ListLoadBalancers().
		Return([]models.LoadBalancerSummary{}, nil)

	c := testutils.NewTestContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadLoadBalancer(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	base.Client.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(&models.LoadBalancer{}, nil)

	c := testutils.NewTestContext(t, []string{"lb_name"}, nil)
	if err := command.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadLoadBalancerInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.read(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

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
			Certificate: "arn:aws:iam::12345:server-certificate/crt_name",
			Expected: models.Port{
				HostPort:       80,
				ContainerPort:  80,
				Protocol:       "https",
				CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name",
			},
		},
		{
			Target:      "80:80/https",
			Certificate: "arn:aws:acm::12345:certificate/crt_name",
			Expected: models.Port{
				HostPort:       80,
				ContainerPort:  80,
				Protocol:       "https",
				CertificateARN: "arn:aws:acm::12345:certificate/crt_name",
			},
		},
	}

	for _, c := range cases {
		result, err := parsePort(c.Target, c.Certificate)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, *result, c.Expected)
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

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := parsePort(c, ""); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestValidateTarget(t *testing.T) {
	cases := []string{
		"TCP:80",
		"HTTP:80/ping/this/path",
		"HTTPS:443/ping/this/path",
	}

	for _, target := range cases {
		t.Run(target, func(t *testing.T) {
			if err := validateTarget(target); err != nil {
				t.Fatal("error was not nil!")
			}
		})
	}
}

func TestValidateTargetErrors(t *testing.T) {
	cases := []string{
		"HTTP:80",
		"HTTPS:443",
	}

	for _, target := range cases {
		t.Run(target, func(t *testing.T) {
			if err := validateTarget(target); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}
