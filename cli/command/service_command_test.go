package command

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/urfave/cli"
	"testing"
)

func TestCreateService(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "environment").
		Return([]string{"environmentID"}, nil)

	tc.Resolver.EXPECT().
		Resolve("deploy", "deploy").
		Return([]string{"deployID"}, nil)

	tc.Resolver.EXPECT().
		Resolve("load_balancer", "load_balancer").
		Return([]string{"loadBalancerID"}, nil)

	tc.Client.EXPECT().
		CreateService("name", "environmentID", "deployID", "loadBalancerID").
		Return(&models.Service{}, nil)

	flags := map[string]interface{}{
		"loadbalancer": "load_balancer",
	}

	c := testutils.GetCLIContext(t, []string{"environment", "name", "deploy"}, flags)
	if err := command.Create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateServiceWait(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "environment").
		Return([]string{"environmentID"}, nil)

	tc.Resolver.EXPECT().
		Resolve("deploy", "deploy").
		Return([]string{"deployID"}, nil)

	tc.Resolver.EXPECT().
		Resolve("load_balancer", "load_balancer").
		Return([]string{"loadBalancerID"}, nil)

	tc.Client.EXPECT().
		CreateService("name", "environmentID", "deployID", "loadBalancerID").
		Return(&models.Service{ServiceID: "serviceID"}, nil)

	tc.Client.EXPECT().
		WaitForDeployment("serviceID", testutils.TEST_TIMEOUT).
		Return(&models.Service{}, nil)

	flags := map[string]interface{}{
		"loadbalancer": "load_balancer",
		"wait":         true,
	}

	c := testutils.GetCLIContext(t, []string{"environment", "name", "deploy"}, flags)
	if err := command.Create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateService_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": testutils.GetCLIContext(t, nil, nil),
		"Missing NAME arg":        testutils.GetCLIContext(t, []string{"environment"}, nil),
		"Missing DEPLOY arg":      testutils.GetCLIContext(t, []string{"environment", "name"}, nil),
	}

	for name, c := range contexts {
		if err := command.Create(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestDeleteService(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("service", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		DeleteService("id").
		Return("jobid", nil)

	c := testutils.GetCLIContext(t, []string{"name"}, nil)
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteServiceWait(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("service", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		DeleteService("id").
		Return("jobid", nil)

	tc.Client.EXPECT().
		WaitForJob("jobid", testutils.TEST_TIMEOUT).
		Return(nil)

	c := testutils.GetCLIContext(t, []string{"name"}, map[string]interface{}{"wait": true})
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteService_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Delete(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestUpdateService(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("service", "service").
		Return([]string{"serviceID"}, nil)

	tc.Resolver.EXPECT().
		Resolve("deploy", "deploy").
		Return([]string{"deployID"}, nil)

	tc.Client.EXPECT().
		UpdateService("serviceID", "deployID").
		Return(&models.Service{}, nil)

	c := testutils.GetCLIContext(t, []string{"service", "deploy"}, nil)
	if err := command.Update(c); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateServiceWait(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("service", "service").
		Return([]string{"serviceID"}, nil)

	tc.Resolver.EXPECT().
		Resolve("deploy", "deploy").
		Return([]string{"deployID"}, nil)

	tc.Client.EXPECT().
		UpdateService("serviceID", "deployID").
		Return(&models.Service{}, nil)

	tc.Client.EXPECT().
		WaitForDeployment("serviceID", testutils.TEST_TIMEOUT).
		Return(&models.Service{}, nil)

	c := testutils.GetCLIContext(t, []string{"service", "deploy"}, map[string]interface{}{"wait": true})
	if err := command.Update(c); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateService_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg":   testutils.GetCLIContext(t, nil, nil),
		"Missing DEPLOY arg": testutils.GetCLIContext(t, []string{"name"}, nil),
	}

	for name, c := range contexts {
		if err := command.Update(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestGetService(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("service", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetService("id").
		Return(&models.Service{}, nil)

	c := testutils.GetCLIContext(t, []string{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetService_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Get(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestListServices(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Client.EXPECT().
		ListServices().
		Return([]*models.ServiceSummary{}, nil)

	c := testutils.GetCLIContext(t, nil, nil)
	if err := command.List(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetServiceLogs(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("service", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetServiceLogs("id", 100)

	c := testutils.GetCLIContext(t, []string{"name"}, map[string]interface{}{"tail": 100})
	if err := command.Logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetServiceLogs_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Logs(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestScaleService(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("service", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		ScaleService("id", 2).
		Return(&models.Service{}, nil)

	c := testutils.GetCLIContext(t, []string{"name", "2"}, nil)
	if err := command.Scale(c); err != nil {
		t.Fatal(err)
	}
}

func TestScaleServiceWait(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("service", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		ScaleService("id", 2).
		Return(&models.Service{}, nil)

	tc.Client.EXPECT().
		WaitForDeployment("id", testutils.TEST_TIMEOUT).
		Return(&models.Service{}, nil)

	c := testutils.GetCLIContext(t, []string{"name", "2"}, map[string]interface{}{"wait": true})
	if err := command.Scale(c); err != nil {
		t.Fatal(err)
	}
}

func TestScaleService_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewServiceCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg":      testutils.GetCLIContext(t, nil, nil),
		"Missing COUNT arg":     testutils.GetCLIContext(t, []string{"name"}, nil),
		"Non-integer COUNT arg": testutils.GetCLIContext(t, []string{"name", "3e"}, nil),
	}

	for name, c := range contexts {
		if err := command.Scale(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}
