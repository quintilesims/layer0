package command

import (
	"net/url"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func TestCreateService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	req := models.CreateServiceRequest{
		DeployID:       "dpl_id",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		ServiceName:    "svc_name",
	}

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	base.Client.EXPECT().
		CreateService(req).
		Return("job_id", nil)

	base.Client.EXPECT().
		ReadJob("job_id").
		Return(&models.Job{
			Status: "Completed",
			Result: "svc_id",
		}, nil)

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(&models.Service{}, nil)

	args := Args{"env_name", "svc_name", "dpl_name"}
	flags := Flags{"loadbalancer": "lb_name"}
	c := getCLIContext(t, args, flags)
	if err := serviceCommand.create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateService_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": getCLIContext(t, nil, nil),
		"Missing NAME arg":        getCLIContext(t, Args{"env_name"}, nil),
		"Missing DEPLOY arg":      getCLIContext(t, Args{"env_name", "svc_name"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.create(c); err == nil {
				t.Fatal("%s: error was nil!", name)
			}
		})
	}
}

func TestDeleteService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		DeleteService("svc_id").
		Return("job_id", nil)

	base.Client.EXPECT().
		ReadJob("job_id").
		Return(&models.Job{
			Status: "Completed",
			Result: "svc_id",
		}, nil)

	args := Args{"svc_name"}
	c := getCLIContext(t, args, nil)
	if err := serviceCommand.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.create(c); err == nil {
				t.Fatal("%s: error was nil!", name)
			}
		})
	}
}

func TestListServices(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	base.Client.EXPECT().
		ListServices().
		Return([]*models.ServiceSummary{}, nil)

	c := getCLIContext(t, nil, nil)
	if err := serviceCommand.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestServiceLogs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		ReadServiceLogs("svc_id", url.Values{
			"tail":  []string{"100"},
			"start": []string{"start"},
			"end":   []string{"end"},
		})

	args := Args{"svc_name"}
	flags := Flags{"tail": 100, "start": "start", "end": "end"}
	c := getCLIContext(t, args, flags)
	if err := serviceCommand.logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestServiceLogs_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.logs(c); err == nil {
				t.Fatal("%s: error was nil!", name)
			}
		})
	}
}

func TestReadService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(&models.Service{}, nil)

	args := Args{"svc_name"}
	c := getCLIContext(t, args, nil)
	if err := serviceCommand.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.read(c); err == nil {
				t.Fatal("%s: error was nil!", name)
			}
		})
	}
}

func TestScaleService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	scale := 2
	req := models.UpdateServiceRequest{
		ServiceID: "svc_id",
		Scale:     &scale,
	}

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		UpdateService(req).
		Return("job_id", nil)

	base.Client.EXPECT().
		ReadJob("job_id").
		Return(&models.Job{
			Status: "Completed",
			Result: "svc_id",
		}, nil)

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(&models.Service{}, nil)

	args := Args{"svc_name", "2"}
	c := getCLIContext(t, args, nil)
	if err := serviceCommand.scale(c); err != nil {
		t.Fatal(err)
	}
}

func TestScaleService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg":  getCLIContext(t, nil, nil),
		"Missing COUNT arg": getCLIContext(t, Args{"svc_name"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.scale(c); err == nil {
				t.Fatal("%s: error was nil!", name)
			}
		})
	}
}

func TestUpdateService(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	deployID := "dpl_id"
	req := models.UpdateServiceRequest{
		ServiceID: "svc_id",
		DeployID:  &deployID,
	}

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	base.Client.EXPECT().
		UpdateService(req).
		Return("job_id", nil)

	base.Client.EXPECT().
		ReadJob("job_id").
		Return(&models.Job{
			Status: "Completed",
			Result: "svc_id",
		}, nil)

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(&models.Service{}, nil)

	args := Args{"svc_name", "dpl_name"}
	c := getCLIContext(t, args, nil)
	if err := serviceCommand.update(c); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	serviceCommand := NewServiceCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg":   getCLIContext(t, nil, nil),
		"Missing DEPLOY arg": getCLIContext(t, Args{"svc_name"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.update(c); err == nil {
				t.Fatal("%s: error was nil!", name)
			}
		})
	}
}
