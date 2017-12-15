package command

import (
	"net/url"
	"testing"

	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func TestCreateService(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
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
			Scale:          3,
		}

		base.Client.EXPECT().
			CreateService(req).
			Return("job_id", nil)

		if wait {
			job := &models.Job{
				Status: "Completed",
				Result: "svc_id",
			}

			base.Client.EXPECT().
				ReadJob("job_id").
				Return(job, nil)

			deployments := []models.Deployment{
				{
					DesiredCount: 3,
					RunningCount: 3,
				},
			}

			service := &models.Service{
				Deployments:  deployments,
				DesiredCount: 3,
				RunningCount: 3,
			}

			base.Client.EXPECT().
				ReadService("svc_id").
				Return(service, nil).
				AnyTimes()
		}

		args := []string{"env_name", "svc_name", "dpl_name"}
		flags := map[string]interface{}{"loadbalancer": "lb_name", "scale": 3}
		c := config.NewTestContext(t, args, flags, config.SetNoWait(!wait))

		serviceCommand := NewServiceCommand(base.Command())
		if err := serviceCommand.create(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCreateService_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": config.NewTestContext(t, nil, nil),
		"Missing NAME arg":        config.NewTestContext(t, []string{"env_name"}, nil),
		"Missing DEPLOY arg":      config.NewTestContext(t, []string{"env_name", "svc_name"}, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.create(c); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestDeleteService(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		base.Resolver.EXPECT().
			Resolve("service", "svc_name").
			Return([]string{"svc_id"}, nil)

		base.Client.EXPECT().
			DeleteService("svc_id").
			Return("job_id", nil)

		if wait {
			job := &models.Job{
				Status: "Completed",
				Result: "svc_id",
			}

			base.Client.EXPECT().
				ReadJob("job_id").
				Return(job, nil)
		}

		args := []string{"svc_name"}
		c := config.NewTestContext(t, args, nil, config.SetNoWait(!wait))

		serviceCommand := NewServiceCommand(base.Command())
		if err := serviceCommand.delete(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDeleteService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing NAME arg": config.NewTestContext(t, nil, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.create(c); err == nil {
				t.Fatal("Error was nil!")
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

	c := config.NewTestContext(t, nil, nil)

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

	args := []string{"svc_name"}
	flags := map[string]interface{}{"tail": 100, "start": "start", "end": "end"}
	c := config.NewTestContext(t, args, flags)

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestServiceLogs_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing NAME arg": config.NewTestContext(t, nil, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.logs(c); err == nil {
				t.Fatal("Error was nil!")
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

	args := []string{"svc_name"}
	c := config.NewTestContext(t, args, nil)

	serviceCommand := NewServiceCommand(base.Command())
	if err := serviceCommand.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing NAME arg": config.NewTestContext(t, nil, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.read(c); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestScaleService(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		defer client.SetTimeMultiplier(0)()

		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		base.Resolver.EXPECT().
			Resolve("service", "svc_name").
			Return([]string{"svc_id"}, nil)

		scale := 2
		req := models.UpdateServiceRequest{Scale: &scale}

		base.Client.EXPECT().
			UpdateService("svc_id", req).
			Return("job_id", nil)

		if wait {
			job := &models.Job{
				Status: "Completed",
				Result: "svc_id",
			}

			base.Client.EXPECT().
				ReadJob("job_id").
				Return(job, nil)

			deployments := []models.Deployment{
				{
					DesiredCount: 2,
					RunningCount: 2,
				},
			}

			service := &models.Service{
				Deployments: deployments,
			}

			base.Client.EXPECT().
				ReadService("svc_id").
				Return(service, nil).
				AnyTimes()
		}

		args := []string{"svc_name", "2"}
		c := config.NewTestContext(t, args, nil, config.SetNoWait(!wait))
		serviceCommand := NewServiceCommand(base.Command())
		if err := serviceCommand.scale(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestScaleService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Resolver.EXPECT().
		Resolve("service", "svc_name").
		Return([]string{"svc_id"}, nil).
		AnyTimes()

	contexts := map[string]*cli.Context{
		"Missing NAME arg":      config.NewTestContext(t, nil, nil),
		"Missing COUNT arg":     config.NewTestContext(t, []string{"svc_name"}, nil),
		"Non-integer COUNT arg": config.NewTestContext(t, []string{"svc_name", "string"}, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.scale(c); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestUpdateService(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
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
		req := models.UpdateServiceRequest{DeployID: &deployID}

		base.Client.EXPECT().
			UpdateService("svc_id", req).
			Return("job_id", nil)

		if wait {
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
		}

		args := []string{"svc_name", "dpl_name"}
		c := config.NewTestContext(t, args, nil, config.SetNoWait(!wait))

		serviceCommand := NewServiceCommand(base.Command())
		if err := serviceCommand.update(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestUpdateService_userInputError(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	contexts := map[string]*cli.Context{
		"Missing NAME arg":   config.NewTestContext(t, nil, nil),
		"Missing DEPLOY arg": config.NewTestContext(t, []string{"svc_name"}, nil),
	}

	serviceCommand := NewServiceCommand(base.Command())
	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := serviceCommand.update(c); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
