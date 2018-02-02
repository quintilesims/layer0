package command

import (
	"net/url"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestCreateService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase()).Command()

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

	input := "l0 service create "
	input += "--loadbalancer lb_name "
	input += "--scale 3 "
	input += "env_name svc_name dpl_name"

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestCreateServiceInputErrors(t *testing.T) {
	testInputErrors(t, NewServiceCommand(nil).Command(), map[string]string{
		"Missing ENVIRONMENT arg": "l0 service create",
		"Missing NAME arg":        "l0 service create env",
		"Missing DEPLOY arg":      "l0 service create env name",
	})
}

func TestDeleteService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		DeleteService("svc_id").
		Return(nil)

	if err := testutils.RunApp(command, "l0 service delete svc_name"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteServiceInputErrors(t *testing.T) {
	testInputErrors(t, NewServiceCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 service delete",
	})
}

func TestListServices(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase()).Command()

	base.Client.EXPECT().
		ListServices().
		Return([]models.ServiceSummary{}, nil)

	if err := testutils.RunApp(command, "l0 service list"); err != nil {
		t.Fatal(err)
	}
}

func TestReadService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(&models.Service{}, nil)

	if err := testutils.RunApp(command, "l0 service get svc_name"); err != nil {
		t.Fatal(err)
	}
}

func TestReadServiceInputErrors(t *testing.T) {
	testInputErrors(t, NewServiceCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 service get",
	})
}

func TestReadServiceLogs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase()).Command()

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

	input := "l0 service logs --tail 100 --start start --end end svc_name"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestReadServiceLogsInputErrors(t *testing.T) {
	testInputErrors(t, NewServiceCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 service logs",
	})
}

func TestScaleService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase()).Command()

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

	if err := testutils.RunApp(command, "l0 service scale svc_name 3"); err != nil {
		t.Fatal(err)
	}
}

func TestScaleServiceInputErrors(t *testing.T) {
	testInputErrors(t, NewServiceCommand(nil).Command(), map[string]string{
		"Missing NAME arg":      "l0 service scale",
		"Missing COUNT arg":     "l0 service scale name",
		"Non-integer COUNT arg": "l0 service scale name two",
	})
}

func TestUpdateService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(base.CommandBase()).Command()

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

	if err := testutils.RunApp(command, "l0 service update svc_name dpl_name"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateServiceInputErrors(t *testing.T) {
	testInputErrors(t, NewServiceCommand(nil).Command(), map[string]string{
		"Missing NAME arg":   "l0 service update",
		"Missing DEPLOY arg": "l0 service update name",
	})
}
