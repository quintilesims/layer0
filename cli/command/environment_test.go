package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func TestEnvironmentCreate_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
		"OS missing":       NewContext(t, Args{"env_name"}, nil),
		"Count negative": NewContext(t, Args{"env_name"},
			Flags{
				"os":        "linux",
				"min-count": "-1",
			}),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.create(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestEnvironmentDelete_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.delete(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestEnvironmentRead_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.read(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestEnvironmentSetMinCount_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg":  NewContext(t, nil, nil),
		"Missing COUNT arg": NewContext(t, Args{"env_name"}, nil),
		"Invalid COUNT arg": NewContext(t, Args{"env_name", "not_a_number"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.update(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestEnvironmentLink_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing SOURCE & DESTINATION args": NewContext(t, nil, nil),
		"Missing DESTINATION arg":           NewContext(t, Args{"env_name1"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.link(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestEnvironmentUnlink_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing SOURCE & DESTINATION args": NewContext(t, nil, nil),
		"Missing DESTINATION arg":           NewContext(t, Args{"env_name1"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.unlink(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestEnvironmentLink_duplicateEnvironmentID(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", gomock.Any()).
		Return([]string{"env_id1"}, nil).
		Times(2)

	c := NewContext(t, []string{"env_name1", "env_name1"}, nil)
	if err := command.link(c); err == nil {
		t.Fatal("error was nil!")
	}
}

func TestEnvironmentUnlink_duplicateEnvironmentID(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", gomock.Any()).
		Return([]string{"env_id1"}, nil).
		Times(2)

	c := NewContext(t, []string{"env_name1", "env_name1"}, nil)
	if err := command.unlink(c); err == nil {
		t.Fatal("error was nil!")
	}
}
func TestCreateEnvironment(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		userData := "user_data"
		file, close := createTempFile(t, userData)
		defer close()

		req := models.CreateEnvironmentRequest{
			EnvironmentName:  "env_name",
			InstanceSize:     "m3.large",
			MinClusterCount:  2,
			UserDataTemplate: []byte(userData),
			OperatingSystem:  "linux",
			AMIID:            "ami",
		}

		environment := &models.Environment{}
		job := &models.Job{
			JobID:  "job_id",
			Status: job.Completed.String(),
			Result: "entity_id",
		}

		base.Client.EXPECT().
			CreateEnvironment(req).
			Return(job.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)

			base.Client.EXPECT().
				ReadEnvironment(job.Result).
				Return(environment, nil)
		}

		flags := map[string]interface{}{
			"size":      req.InstanceSize,
			"min-count": req.MinClusterCount,
			"user-data": file.Name(),
			"os":        req.OperatingSystem,
			"ami":       req.AMIID,
		}

		c := NewContext(t, []string{"env_name"}, flags, SetNoWait(!wait))
		if err := command.create(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDeleteEnvironment(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		job := &models.Job{
			JobID:  "job_id",
			Status: job.Completed.String(),
			Result: "entity_id",
		}

		base.Resolver.EXPECT().
			Resolve("environment", "env_name").
			Return([]string{"env_id"}, nil)

		base.Client.EXPECT().
			DeleteEnvironment("env_id").
			Return(job.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)
		}

		c := NewContext(t, []string{"env_name"}, nil, SetNoWait(!wait))
		if err := command.delete(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetEnvironment(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", "env_name*").
		Return([]string{"env_id1", "env_id2"}, nil)

	base.Client.EXPECT().
		ReadEnvironment("env_id1").
		Return(&models.Environment{}, nil)

	base.Client.EXPECT().
		ReadEnvironment("env_id2").
		Return(&models.Environment{}, nil)

	c := NewContext(t, []string{"env_name*"}, nil)
	if err := command.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestListEnvironments(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	base.Client.EXPECT().
		ListEnvironments().
		Return([]*models.EnvironmentSummary{}, nil)

	c := NewContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentSetMinCount(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		minCount := 2
		job := &models.Job{
			JobID:  "job_id",
			Status: job.Completed.String(),
			Result: "entity_id",
		}

		base.Resolver.EXPECT().
			Resolve("environment", "env_name").
			Return([]string{"env_id"}, nil)

		req := models.UpdateEnvironmentRequest{
			EnvironmentID:   "env_id",
			MinClusterCount: &minCount,
		}

		base.Client.EXPECT().
			UpdateEnvironment(req).
			Return(job.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)

			base.Client.EXPECT().
				ReadEnvironment(job.Result).
				Return(&models.Environment{}, nil)
		}

		c := NewContext(t, []string{"env_name", "2"}, nil, SetNoWait(!wait))
		if err := command.update(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestEnvironmentLink(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		job := &models.Job{
			JobID:  "job_id",
			Status: job.Completed.String(),
			Result: "entity_id",
		}

		base.Resolver.EXPECT().
			Resolve("environment", "env_name1").
			Return([]string{"env_id1"}, nil)

		base.Resolver.EXPECT().
			Resolve("environment", "env_name2").
			Return([]string{"env_id2"}, nil)

		req := models.CreateEnvironmentLinkRequest{
			SourceEnvironmentID: "env_id1",
			DestEnvironmentID:   "env_id2",
		}
		base.Client.EXPECT().
			CreateLink(req).
			Return(job.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)
		}

		c := NewContext(t, []string{"env_name1", "env_name2"}, nil, SetNoWait(!wait))
		if err := command.link(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestEnvironmentUnlink(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		job := &models.Job{
			JobID:  "job_id",
			Status: job.Completed.String(),
			Result: "entity_id",
		}

		base.Resolver.EXPECT().
			Resolve("environment", "env_name1").
			Return([]string{"env_id1"}, nil)

		base.Resolver.EXPECT().
			Resolve("environment", "env_name2").
			Return([]string{"env_id2"}, nil)

		req := models.DeleteEnvironmentLinkRequest{
			SourceEnvironmentID: "env_id1",
			DestEnvironmentID:   "env_id2",
		}
		base.Client.EXPECT().
			DeleteLink(req).
			Return(job.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)
		}

		c := NewContext(t, []string{"env_name1", "env_name2"}, nil, SetNoWait(!wait))
		if err := command.unlink(c); err != nil {
			t.Fatal(err)
		}
	})
}
