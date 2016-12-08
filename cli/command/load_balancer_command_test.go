package command

import (
	"github.com/urfave/cli"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestParsePort(t *testing.T) {
	cases := map[string]*models.Port{
		"80:80/tcp": &models.Port{
			HostPort:      80,
			ContainerPort: 80,
			Protocol:      "tcp",
			CertificateID: "",
		},
		"80:80/http": &models.Port{
			HostPort:      80,
			ContainerPort: 80,
			Protocol:      "http",
			CertificateID: "",
		},
		"8080:80/http": &models.Port{
			HostPort:      8080,
			ContainerPort: 80,
			Protocol:      "http",
			CertificateID: "",
		},
		"443:80/https": &models.Port{
			HostPort:      443,
			ContainerPort: 80,
			Protocol:      "https",
			CertificateID: "certid",
		},
	}

	for input, expected := range cases {
		model, err := parsePort(input, "certid")
		if err != nil {
			t.Fatal(err)
		}

		testutils.AssertEqual(t, model, expected)
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
		HostPort:      443,
		ContainerPort: 80,
		Protocol:      "https",
		CertificateID: "certid",
	}

	tc.Client.EXPECT().
		UpdateLoadBalancer("id", []models.Port{port}).
		Return(&models.LoadBalancer{}, nil)

	flags := Flags{"certificate": "certid"}
	c := getCLIContext(t, Args{"name", "443:80/https"}, flags)
	if err := command.AddPort(c); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerAddPort_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
		"Missing PORT arg": getCLIContext(t, Args{"name"}, nil),
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

	ports := []models.Port{
		{
			HostPort:      443,
			ContainerPort: 80,
			Protocol:      "https",
			CertificateID: "certid",
		},
		{
			HostPort:      8000,
			ContainerPort: 8000,
			Protocol:      "http",
			CertificateID: "",
		},
	}

	tc.Client.EXPECT().
		CreateLoadBalancer("name", "environmentID", ports, false).
		Return(&models.LoadBalancer{}, nil)

	flags := Flags{
		"port":        []string{"443:80/https", "8000:8000/http"},
		"certificate": "certid",
		"private":     true,
	}

	c := getCLIContext(t, Args{"environment", "name"}, flags)
	if err := command.Create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateLoadBalancer_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": getCLIContext(t, nil, nil),
		"Missing NAME arg":        getCLIContext(t, Args{"environment"}, nil),
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

	c := getCLIContext(t, Args{"name"}, nil)
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
		WaitForJob("jobid", TEST_TIMEOUT).
		Return(nil)

	c := getCLIContext(t, Args{"name"}, Flags{"wait": true})
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteLoadBalancer_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
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

	c := getCLIContext(t, Args{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetLoadBalancer_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
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
		Return([]*models.LoadBalancer{}, nil)

	c := getCLIContext(t, nil, nil)
	if err := command.List(c); err != nil {
		t.Fatal(err)
	}
}
