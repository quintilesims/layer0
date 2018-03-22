package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/stretchr/testify/assert"
)

func TestLoadBalancerAddPort(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.CommandBase()).Command()

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

	input := "l0 loadbalancer addport "
	input += "--certificate arn:aws:iam::12345:server-certificate/crt_name "
	input += "lb_name 443:80/https"

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerAddPortInputErrors(t *testing.T) {
	testInputErrors(t, NewLoadBalancerCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 loadbalancer addport",
		"Missing PORT arg": "l0 loadbalancer addport name",
	})
}

func TestCreateLoadBalancer(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		LoadBalancerType: config.DefaultLoadBalancerType,
		EnvironmentID:    "env_id",
		IsPublic:         false,
		Ports: []models.Port{
			{HostPort: 443, ContainerPort: 80, Protocol: "https", CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name"},
			{HostPort: 22, ContainerPort: 22, Protocol: "tcp", CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name"},
		},
		HealthCheck: models.HealthCheck{
			Target:             "tcp:80",
			Path:               "/",
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

	input := "l0 loadbalancer create "
	input += "--private "
	input += "--type " + string(config.DefaultLoadBalancerType) + " "
	input += "--certificate arn:aws:iam::12345:server-certificate/crt_name "
	input += "--port 443:80/https "
	input += "--port 22:22/tcp "
	input += "--healthcheck-target tcp:80 "
	input += "--healthcheck-interval 5 "
	input += "--healthcheck-timeout 6 "
	input += "--healthcheck-healthy-threshold 7 "
	input += "--healthcheck-unhealthy-threshold 8 "
	input += "env_name lb_name"

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestCreateLoadBalancerInputErrors(t *testing.T) {
	testInputErrors(t, NewLoadBalancerCommand(nil).Command(), map[string]string{
		"Missing ENVIRONMENT arg": "l0 loadbalancer create",
		"Missing NAME arg":        "l0 loadbalancer create env",
	})
}

func TestDeleteLoadBalancer(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	base.Client.EXPECT().
		DeleteLoadBalancer("lb_id").
		Return(nil)

	input := "l0 loadbalancer delete lb_name"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteLoadBalancerInputErrors(t *testing.T) {
	testInputErrors(t, NewLoadBalancerCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 loadbalancer delete",
	})
}

func TestLoadBalancerDropPort(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.CommandBase()).Command()

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

	input := "l0 loadbalancer dropport lb_name 443"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerDropPortInputErrors(t *testing.T) {
	testInputErrors(t, NewLoadBalancerCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 loadbalancer dropport",
		"Missing PORT arg": "l0 loadbalancer dropport name",
	})
}

func TestListLoadBalancers(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.CommandBase()).Command()

	base.Client.EXPECT().
		ListLoadBalancers().
		Return([]models.LoadBalancerSummary{}, nil)

	input := "l0 loadbalancer list"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestReadLoadBalancer(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	base.Client.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(&models.LoadBalancer{}, nil)

	input := "l0 loadbalancer get lb_name"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestReadLoadBalancerInputErrors(t *testing.T) {
	testInputErrors(t, NewLoadBalancerCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 loadbalancer get",
	})
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
		hc := models.HealthCheck{
			Target:             target,
			Interval:           1,
			Timeout:            1,
			HealthyThreshold:   1,
			UnhealthyThreshold: 1,
		}
		t.Run(target, func(t *testing.T) {
			if err := hc.Validate(); err != nil {
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
		hc := models.HealthCheck{
			Target:             target,
			Interval:           1,
			Timeout:            1,
			HealthyThreshold:   1,
			UnhealthyThreshold: 1,
		}
		t.Run(target, func(t *testing.T) {
			if err := hc.Validate(); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	cases := []string{
		"/",
		"/ping/this/path",
	}

	for _, path := range cases {
		hc := models.HealthCheck{
			Path:               path,
			Interval:           1,
			Timeout:            1,
			HealthyThreshold:   1,
			UnhealthyThreshold: 1,
		}
		t.Run(path, func(t *testing.T) {
			if err := hc.Validate(); err != nil {
				t.Fatal("error was not nil!")
			}
		})
	}
}

func TestValidatePathErrors(t *testing.T) {
	cases := []string{
		"ping/this/path",
		"TCP:80",
		"HTTP:80/ping/this/path",
	}

	for _, path := range cases {
		hc := models.HealthCheck{
			Path:               path,
			Interval:           1,
			Timeout:            1,
			HealthyThreshold:   1,
			UnhealthyThreshold: 1,
		}
		t.Run(path, func(t *testing.T) {
			if err := hc.Validate(); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}
