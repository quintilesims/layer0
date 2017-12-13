package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestAddPortToLoadBalancer(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		base.Resolver.EXPECT().
			Resolve("load_balancer", "lb_name").
			Return([]string{"lb_id"}, nil)

		ports := []models.Port{
			models.Port{
				ContainerPort: 80,
				HostPort:      80,
				Protocol:      "tcp",
			},
		}

		loadBalancer := &models.LoadBalancer{
			LoadBalancerID:   "lb_id",
			LoadBalancerName: "lb_name",
			Ports:            ports,
		}

		base.Client.EXPECT().
			ReadLoadBalancer("lb_id").
			Return(loadBalancer, nil)

		ports = append(ports, models.Port{
			ContainerPort: 81,
			HostPort:      81,
			Protocol:      "tcp",
		})

		req := models.UpdateLoadBalancerRequest{Ports: &ports}

		base.Client.EXPECT().
			UpdateLoadBalancer("lb_id", req).
			Return("jid", nil)

		job := &models.Job{
			Status: "Completed",
			Result: "lb_id",
		}

		if wait {
			base.Client.EXPECT().
				ReadJob("jid").
				Return(job, nil)

			base.Client.EXPECT().
				ReadLoadBalancer("lb_id").
				Return(&models.LoadBalancer{}, nil)
		}

		args := Args{"lb_name", "81:81/tcp"}
		c := NewContext(t, args, nil, SetNoWait(!wait))

		loadBalancerCommand := NewLoadBalancerCommand(base.Command())
		if err := loadBalancerCommand.addport(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestAddPortToLoadBalancer_UserInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	loadBalancerCommand := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
		"Missing PORT arg": NewContext(t, Args{"name"}, nil),
	}

	for name, c := range contexts {
		if err := loadBalancerCommand.addport(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestCreateLoadBalancer(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		base.Resolver.EXPECT().
			Resolve("environment", "ename").
			Return([]string{"eid"}, nil)

		ports := []models.Port{
			models.Port{
				ContainerPort: 80,
				HostPort:      80,
				Protocol:      "tcp",
			},
		}

		healthCheck := config.DefaultLoadBalancerHealthCheck

		req := models.CreateLoadBalancerRequest{
			LoadBalancerName: "lb_name",
			EnvironmentID:    "eid",
			IsPublic:         true,
			Ports:            ports,
			HealthCheck:      healthCheck,
		}

		base.Client.EXPECT().
			CreateLoadBalancer(req).
			Return("jid", nil)

		job := &models.Job{
			Status: "Completed",
			Result: "lb_id",
		}

		if wait {
			base.Client.EXPECT().
				ReadJob("jid").
				Return(job, nil)

			base.Client.EXPECT().
				ReadLoadBalancer("lb_id").
				Return(&models.LoadBalancer{}, nil)
		}

		args := Args{"ename", "lb_name"}
		flags := Flags{
			"healthcheck-target":              healthCheck.Target,
			"healthcheck-interval":            healthCheck.Interval,
			"healthcheck-timeout":             healthCheck.Timeout,
			"healthcheck-healthy-threshold":   healthCheck.HealthyThreshold,
			"healthcheck-unhealthy-threshold": healthCheck.UnhealthyThreshold,
			"port": []string{"80:80/tcp"},
		}

		c := NewContext(t, args, flags, SetNoWait(!wait))

		loadBalancerCommand := NewLoadBalancerCommand(base.Command())
		if err := loadBalancerCommand.create(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCreateLoadBalancer_UserInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	loadBalancerCommand := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": NewContext(t, nil, nil),
		"Missing NAME arg":        NewContext(t, Args{"environment"}, nil),
	}

	for name, c := range contexts {
		if err := loadBalancerCommand.create(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestDeleteLoadBalancer(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		base.Resolver.EXPECT().
			Resolve("load_balancer", "lb_name").
			Return([]string{"lb_id"}, nil)

		base.Client.EXPECT().
			DeleteLoadBalancer("lb_id").
			Return("jid", nil)

		job := &models.Job{
			Status: "Completed",
			Result: "lb_id",
		}

		if wait {
			base.Client.EXPECT().
				ReadJob("jid").
				Return(job, nil)
		}

		args := Args{"lb_name"}
		c := NewContext(t, args, nil, SetNoWait(!wait))

		loadBalancerCommand := NewLoadBalancerCommand(base.Command())
		if err := loadBalancerCommand.delete(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDeleteLoadBalancer_UserInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	loadBalancerCommand := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := loadBalancerCommand.delete(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestDropPortFromLoadBalancer(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		base.Resolver.EXPECT().
			Resolve("load_balancer", "lb_name").
			Return([]string{"lb_id"}, nil)

		ports := []models.Port{
			models.Port{
				ContainerPort: 80,
				HostPort:      80,
				Protocol:      "tcp",
			},
			models.Port{
				ContainerPort: 81,
				HostPort:      81,
				Protocol:      "tcp",
			},
		}

		loadBalancer := &models.LoadBalancer{
			LoadBalancerID:   "lb_id",
			LoadBalancerName: "lb_name",
			Ports:            ports,
		}

		base.Client.EXPECT().
			ReadLoadBalancer("lb_id").
			Return(loadBalancer, nil)

		ports = ports[:1]

		req := models.UpdateLoadBalancerRequest{Ports: &ports}

		base.Client.EXPECT().
			UpdateLoadBalancer("lb_id", req).
			Return("jid", nil)

		job := &models.Job{
			Status: "Completed",
			Result: "lb_id",
		}

		if wait {
			base.Client.EXPECT().
				ReadJob("jid").
				Return(job, nil)

			base.Client.EXPECT().
				ReadLoadBalancer("lb_id").
				Return(&models.LoadBalancer{}, nil)
		}

		args := Args{"lb_name", "81"}
		c := NewContext(t, args, nil, SetNoWait(!wait))

		loadBalancerCommand := NewLoadBalancerCommand(base.Command())
		if err := loadBalancerCommand.dropport(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDropPortFromLoadBalancer_UserInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	loadBalancerCommand := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
		"Missing PORT arg": NewContext(t, Args{"name"}, nil),
	}

	for name, c := range contexts {
		if err := loadBalancerCommand.dropport(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestDisplayHealthCheck(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	loadBalancer := &models.LoadBalancer{
		LoadBalancerID:   "lb_id",
		LoadBalancerName: "lb_name",
	}

	base.Client.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(loadBalancer, nil)

	args := Args{"lb_name"}
	c := NewContext(t, args, nil)

	loadBalancerCommand := NewLoadBalancerCommand(base.Command())
	if err := loadBalancerCommand.healthcheck(c); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateHealthCheck(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		base.Resolver.EXPECT().
			Resolve("load_balancer", "lb_name").
			Return([]string{"lb_id"}, nil)

		healthCheck := models.HealthCheck{Target: "TCP:80"}
		loadBalancer := &models.LoadBalancer{
			HealthCheck:      healthCheck,
			LoadBalancerID:   "lb_id",
			LoadBalancerName: "lb_name",
		}

		base.Client.EXPECT().
			ReadLoadBalancer("lb_id").
			Return(loadBalancer, nil)

		healthCheck = models.HealthCheck{Target: "TCP:81"}
		req := models.UpdateLoadBalancerRequest{HealthCheck: &healthCheck}

		base.Client.EXPECT().
			UpdateLoadBalancer("lb_id", req).
			Return("jid", nil)

		job := &models.Job{
			Status: "Completed",
			Result: "lb_id",
		}

		if wait {
			base.Client.EXPECT().
				ReadJob("jid").
				Return(job, nil)

			base.Client.EXPECT().
				ReadLoadBalancer("lb_id").
				Return(&models.LoadBalancer{}, nil)
		}

		args := Args{"lb_name"}
		flags := Flags{"healthcheck-target": "TCP:81"}
		c := NewContext(t, args, flags, SetNoWait(!wait))

		loadBalancerCommand := NewLoadBalancerCommand(base.Command())
		if err := loadBalancerCommand.healthcheck(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestUpdateHealthCheck_UserInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{"Missing NAME arg": NewContext(t, nil, nil)}

	loadBalancerCommand := NewLoadBalancerCommand(base.Command())
	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := loadBalancerCommand.healthcheck(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestListLoadBalancers(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Client.EXPECT().
		ListLoadBalancers()

	c := NewContext(t, nil, nil)

	loadBalancerCommand := NewLoadBalancerCommand(base.Command())
	if err := loadBalancerCommand.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadLoadBalancer(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	loadBalancerIDs := []string{"lb_id1", "lb_id2"}

	base.Resolver.EXPECT().
		Resolve("load_balancer", "*").
		Return(loadBalancerIDs, nil)

	for _, id := range loadBalancerIDs {
		base.Client.EXPECT().
			ReadLoadBalancer(id).
			Return(&models.LoadBalancer{}, nil)
	}

	args := Args{"*"}
	c := NewContext(t, args, nil)

	loadBalancerCommand := NewLoadBalancerCommand(base.Command())
	if err := loadBalancerCommand.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadLoadBalancer_UserInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	loadBalancerCommand := NewLoadBalancerCommand(base.Command())
	for name, c := range contexts {
		if err := loadBalancerCommand.read(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestParsePort(t *testing.T) {
	cases := map[string]*models.Port{
		"80:80/tcp": {
			HostPort:        80,
			ContainerPort:   80,
			Protocol:        "tcp",
			CertificateName: "",
		},
		"80:80/http": {
			HostPort:        80,
			ContainerPort:   80,
			Protocol:        "http",
			CertificateName: "",
		},
		"8080:80/http": {
			HostPort:        8080,
			ContainerPort:   80,
			Protocol:        "http",
			CertificateName: "",
		},
		"443:80/https": {
			HostPort:        443,
			ContainerPort:   80,
			Protocol:        "https",
			CertificateName: "cert_name",
		},
	}

	for input, expected := range cases {
		model, err := parsePort(input, expected.CertificateName)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, model)
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
				t.Fatalf("%s: error was nil!", name)
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
				t.Fatalf("%s: error was not nil!", target)
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
				t.Fatalf("%s: error was nil!", target)
			}
		})
	}
}
