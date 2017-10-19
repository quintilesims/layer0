package command

import (
	"net/url"
	"testing"

	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func TestCreateService(t *testing.T) {
	defer client.SetTimeMultiplier(0)()

	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	req := models.CreateServiceRequest{
		DeployID:       "dpl_id",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		ServiceName:    "svc_name",
	}

	base.Client.EXPECT().
		CreateService(req).
		Return("job_id", nil)

	job := &models.Job{
		Status: "Completed",
		Result: "svc_id",
	}

	base.Client.EXPECT().
		ReadJob("job_id").
		Return(job, nil)

	deployments := []models.Deployment{
		{
			DesiredCount: 1,
			RunningCount: 1,
		},
	}

	service := &models.Service{
		Deployments:  deployments,
		DesiredCount: 1,
		RunningCount: 1,
	}

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(service, nil).
		AnyTimes()

	args := Args{"env_name", "svc_name", "dpl_name"}
	flags := Flags{"loadbalancer": "lb_name"}
	c := NewContext(t, args, flags)

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateService_noWait(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("load_balancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	req := models.CreateServiceRequest{
		DeployID:       "dpl_id",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		ServiceName:    "svc_name",
	}

	base.Client.EXPECT().
		CreateService(req).
		Return("job_id", nil)

	args := Args{"env_name", "svc_name", "dpl_name"}
	flags := Flags{"loadbalancer": "lb_name"}
	c := NewContext(t, args, flags, SetNoWait(true))

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateService_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": NewContext(t, nil, nil),
		"Missing NAME arg":        NewContext(t, Args{"env_name"}, nil),
		"Missing DEPLOY arg":      NewContext(t, Args{"env_name", "svc_name"}, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
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

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		DeleteService("svc_id").
		Return("job_id", nil)

	job := &models.Job{
		Status: "Completed",
		Result: "svc_id",
	}

	base.Client.EXPECT().
		ReadJob("job_id").
		Return(job, nil)

	args := Args{"svc_name"}
	c := NewContext(t, args, nil)

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteService_noWait(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		DeleteService("svc_id").
		Return("job_id", nil)

	args := Args{"svc_name"}
	c := NewContext(t, args, nil, SetNoWait(true))

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
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

	base.Client.EXPECT().
		ListServices().
		Return([]*models.ServiceSummary{}, nil)

	c := NewContext(t, nil, nil)

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestServiceLogs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

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
		Return([]*models.LogFile{}, nil)

	args := Args{"svc_name"}
	flags := Flags{"tail": 100, "start": "start", "end": "end"}
	c := NewContext(t, args, flags)

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestServiceLogs_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
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

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(&models.Service{}, nil)

	args := Args{"svc_name"}
	c := NewContext(t, args, nil)

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.read(c); err == nil {
				t.Fatal("%s: error was nil!", name)
			}
		})
	}
}

func TestScaleService(t *testing.T) {
	defer client.SetTimeMultiplier(0)()

	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	scale := 2
	req := models.UpdateServiceRequest{
		ServiceID: "svc_id",
		Scale:     &scale,
	}

	base.Client.EXPECT().
		UpdateService(req).
		Return("job_id", nil)

	job := &models.Job{
		Status: "Completed",
		Result: "svc_id",
	}

	base.Client.EXPECT().
		ReadJob("job_id").
		Return(job, nil)

	deployments := []models.Deployment{
		{
			DesiredCount: 1,
			RunningCount: 1,
		},
	}

	service := &models.Service{
		Deployments:  deployments,
		DesiredCount: 1,
		RunningCount: 1,
	}

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(service, nil).
		AnyTimes()

	args := Args{"svc_name", "2"}
	c := NewContext(t, args, nil)
	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.scale(c); err != nil {
		t.Fatal(err)
	}
}

func TestScaleService_noWait(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	scale := 2
	req := models.UpdateServiceRequest{
		ServiceID: "svc_id",
		Scale:     &scale,
	}

	base.Client.EXPECT().
		UpdateService(req).
		Return("job_id", nil)

	args := Args{"svc_name", "2"}
	c := NewContext(t, args, nil, SetNoWait(true))

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.scale(c); err != nil {
		t.Fatal(err)
	}
}

func TestScaleService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil).
		AnyTimes()

	contexts := map[string]*cli.Context{
		"Missing NAME arg":      NewContext(t, nil, nil),
		"Missing COUNT arg":     NewContext(t, Args{"svc_name"}, nil),
		"Non-integer COUNT arg": NewContext(t, Args{"svc_name", "string"}, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
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

	defer client.SetTimeMultiplier(0)()

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	deployID := "dpl_id"
	req := models.UpdateServiceRequest{
		ServiceID: "svc_id",
		DeployID:  &deployID,
	}

	base.Client.EXPECT().
		UpdateService(req).
		Return("job_id", nil)

	job := &models.Job{
		Status: "Completed",
		Result: "svc_id",
	}

	base.Client.EXPECT().
		ReadJob("job_id").
		Return(job, nil)

	deployments := []models.Deployment{
		{
			DesiredCount: 1,
			RunningCount: 1,
		},
	}

	service := &models.Service{
		Deployments:  deployments,
		DesiredCount: 1,
		RunningCount: 1,
	}

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(service, nil).
		AnyTimes()

	args := Args{"svc_name", "dpl_name"}
	c := NewContext(t, args, nil)

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.update(c); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateService_noWait(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	deployID := "dpl_id"
	req := models.UpdateServiceRequest{
		ServiceID: "svc_id",
		DeployID:  &deployID,
	}

	base.Client.EXPECT().
		UpdateService(req).
		Return("job_id", nil)

	args := Args{"svc_name", "dpl_name"}
	c := NewContext(t, args, nil, SetNoWait(true))

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.update(c); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing NAME arg":   NewContext(t, nil, nil),
		"Missing DEPLOY arg": NewContext(t, Args{"svc_name"}, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.update(c); err == nil {
				t.Fatal("%s: error was nil!", name)
			}
		})
	}
}
