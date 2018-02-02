package command

import (
	"net/url"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/urfave/cli"
)

func TestCreateService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	req := models.CreateServiceRequest{
		ServiceName:    "svc_name",
		EnvironmentID:  "env_id",
		DeployID:       "dpl_id",
		LoadBalancerID: "lb_id",
		Scale:          3,
	}

	base.Client.EXPECT().
		CreateService(req).
		Return("svc_id", nil)

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(&models.Service{}, nil)

	flags := map[string]interface{}{
		"loadbalancer": "lb_name",
		"scale":        3,
	}

	c := testutils.NewTestContext(t, []string{"env_name", "svc_name", "dpl_name"}, flags)
	if err := command.create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateServiceInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	cases := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": testutils.NewTestContext(t, nil, nil),
		"Missing NAME arg":        testutils.NewTestContext(t, []string{"env_name"}, nil),
		"Missing DEPLOY arg":      testutils.NewTestContext(t, []string{"env_name", "svc_name"}, nil),
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := command.create(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestDeleteService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		DeleteService("svc_id").
		Return(nil)

	c := testutils.NewTestContext(t, []string{"svc_name"}, nil)
	if err := command.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteServiceInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	cases := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := command.delete(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestListServices(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	base.Client.EXPECT().
		ListServices().
		Return([]models.ServiceSummary{}, nil)

	c := testutils.NewTestContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(&models.Service{}, nil)

	c := testutils.NewTestContext(t, []string{"svc_name"}, nil)
	if err := command.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadServiceInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	cases := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := command.read(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestReadServiceLogs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	query := url.Values{
		"tail":  []string{"100"},
		"start": []string{"start"},
		"end":   []string{"end"},
	}

	base.Client.EXPECT().
		ReadServiceLogs("svc_id", query).
		Return([]models.LogFile{}, nil)

	flags := map[string]interface{}{
		"tail":  100,
		"start": "start",
		"end":   "end",
	}

	c := testutils.NewTestContext(t, []string{"svc_name"}, flags)
	if err := command.logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadServiceLogsInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	cases := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := command.logs(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestScaleService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	scale := 3
	req := models.UpdateServiceRequest{
		Scale: &scale,
	}

	base.Client.EXPECT().
		UpdateService("svc_id", req).
		Return(nil)

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(&models.Service{}, nil)

	c := testutils.NewTestContext(t, []string{"svc_name", "3"}, nil)
	if err := command.scale(c); err != nil {
		t.Fatal(err)
	}
}

func TestScaleServiceInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	cases := map[string]*cli.Context{
		"Missing NAME arg":      testutils.NewTestContext(t, nil, nil),
		"Missing COUNT arg":     testutils.NewTestContext(t, []string{"svc_name"}, nil),
		"Non-integer COUNT arg": testutils.NewTestContext(t, []string{"svc_name", "two"}, nil),
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := command.scale(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestUpdateService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	deployID := "dpl_id"
	req := models.UpdateServiceRequest{
		DeployID: &deployID,
	}

	base.Client.EXPECT().
		UpdateService("svc_id", req).
		Return(nil)

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(&models.Service{}, nil)

	c := testutils.NewTestContext(t, []string{"svc_name", "dpl_name"}, nil)
	if err := command.update(c); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateServiceInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase())

	cases := map[string]*cli.Context{
		"Missing NAME arg":   testutils.NewTestContext(t, nil, nil),
		"Missing DEPLOY arg": testutils.NewTestContext(t, []string{"svc_name"}, nil),
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := command.update(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}
