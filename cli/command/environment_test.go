package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func TestEnvironmentCommand_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	testCases := []struct {
		name    string
		command func(*cli.Context) error
		args    []string
	}{
		{
			name:    "create",
			command: command.create,
		},
		{
			name:    "update",
			command: command.update,
		},
		{
			name:    "setMinCount",
			command: command.update,
			args:    []string{"env_name", "1w"},
		},
		{
			name:    "read",
			command: command.read,
		},
		{
			name:    "delete",
			command: command.delete,
		},
		{
			name:    "link",
			command: command.link,
		},
		{
			name:    "unlink",
			command: command.unlink,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := getCLIContext(t, tc.args, nil)

			if err := tc.command(c); err == nil {
				t.Fatalf("%s: error was nil!", tc.name)
			}
		})
	}
}

func TestCreateEnvironment(t *testing.T) {
	testNoWaitCaseHelper(t, func(t *testing.T, otherFlags map[string]interface{}, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		userData := "user_data"
		file, close := tempFile(t, userData)
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
		for k, v := range otherFlags {
			flags[k] = v
		}

		c := getCLIContext(t, []string{"env_name"}, flags)
		if err := command.create(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDeleteEnvironment(t *testing.T) {
	testNoWaitCaseHelper(t, func(t *testing.T, flags map[string]interface{}, wait bool) {
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

		c := getCLIContext(t, []string{"env_name"}, flags)
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
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Client.EXPECT().
		ReadEnvironment("env_id").
		Return(&models.Environment{}, nil)

	c := getCLIContext(t, []string{"env_name"}, nil)
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

	c := getCLIContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentSetMinCount(t *testing.T) {
	testNoWaitCaseHelper(t, func(t *testing.T, flags map[string]interface{}, wait bool) {
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
			UpdateEnvironment(gomock.Any()).
			Return(job.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)

			base.Client.EXPECT().
				ReadEnvironment(job.Result).
				Return(&models.Environment{}, nil)
		}

		c := getCLIContext(t, []string{"env_name", "2"}, flags)
		if err := command.update(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestEnvironmentLink(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", "env_name1").
		Return([]string{"env_id1"}, nil)

	base.Resolver.EXPECT().
		Resolve("environment", "env_name2").
		Return([]string{"env_id2"}, nil)

	base.Client.EXPECT().
		CreateLink("env_id1", "env_id2").
		Return(nil)

	c := getCLIContext(t, []string{"env_name1", "env_name2"}, nil)
	if err := command.link(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentUnlink(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", "env_name1").
		Return([]string{"env_id1"}, nil)

	base.Resolver.EXPECT().
		Resolve("environment", "env_name2").
		Return([]string{"env_id2"}, nil)

	base.Client.EXPECT().
		DeleteLink("env_id1", "env_id2").
		Return(nil)

	c := getCLIContext(t, []string{"env_name1", "env_name2"}, nil)
	if err := command.unlink(c); err != nil {
		t.Fatal(err)
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

	c := getCLIContext(t, []string{"env_name1", "env_name1"}, nil)
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

	c := getCLIContext(t, []string{"env_name1", "env_name1"}, nil)
	if err := command.unlink(c); err == nil {
		t.Fatal("error was nil!")
	}
}
