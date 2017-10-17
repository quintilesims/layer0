package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/cli/resolver/mock_resolver"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func TestEnvironmentCommand_userInputErrors(t *testing.T) {
	_, _, command, _ := initEnvCommandTest(t)

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
			args:    []string{"NAME", "1w"},
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
		client, _, command, ctrl := initEnvCommandTest(t)
		defer ctrl.Finish()

		userData := "user_data"
		file, close := tempFile(t, userData)
		defer close()

		req := models.CreateEnvironmentRequest{
			EnvironmentName:  "name",
			InstanceSize:     "m3.large",
			MinClusterCount:  2,
			UserDataTemplate: []byte(userData),
			OperatingSystem:  "linux",
			AMIID:            "ami",
		}

		environment := &models.Environment{}
		job := &models.Job{
			JobID:  "job-id",
			Status: job.Completed.String(),
			Result: "entity-id",
		}

		client.EXPECT().
			CreateEnvironment(req).
			Return(job.JobID, nil)

		if wait {
			client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)

			client.EXPECT().
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

		c := getCLIContext(t, []string{"name"}, flags)
		if err := command.create(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDeleteEnvironment(t *testing.T) {
	testNoWaitCaseHelper(t, func(t *testing.T, flags map[string]interface{}, wait bool) {
		client, resolver, command, ctrl := initEnvCommandTest(t)
		defer ctrl.Finish()

		job := &models.Job{
			JobID:  "job-id",
			Status: job.Completed.String(),
			Result: "entity-id",
		}

		resolver.EXPECT().
			Resolve("environment", "name").
			Return([]string{"id"}, nil)

		client.EXPECT().
			DeleteEnvironment("id").
			Return(job.JobID, nil)

		if wait {
			client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)
		}

		c := getCLIContext(t, []string{"name"}, flags)
		if err := command.delete(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetEnvironment(t *testing.T) {
	client, resolver, command, ctrl := initEnvCommandTest(t)
	defer ctrl.Finish()

	resolver.EXPECT().
		Resolve("environment", "name").
		Return([]string{"id"}, nil)

	client.EXPECT().
		ReadEnvironment("id").
		Return(&models.Environment{}, nil)

	c := getCLIContext(t, []string{"name"}, nil)
	if err := command.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestListEnvironments(t *testing.T) {
	client, _, command, ctrl := initEnvCommandTest(t)
	defer ctrl.Finish()

	client.EXPECT().
		ListEnvironments().
		Return([]*models.EnvironmentSummary{}, nil)

	c := getCLIContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentSetMinCount(t *testing.T) {
	testNoWaitCaseHelper(t, func(t *testing.T, flags map[string]interface{}, wait bool) {
		client, resolver, command, ctrl := initEnvCommandTest(t)
		defer ctrl.Finish()

		job := &models.Job{
			JobID:  "job-id",
			Status: job.Completed.String(),
			Result: "entity-id",
		}

		resolver.EXPECT().
			Resolve("environment", "name").
			Return([]string{"id"}, nil)

		client.EXPECT().
			UpdateEnvironment(gomock.Any()).
			Return(job.JobID, nil)

		if wait {
			client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)

			client.EXPECT().
				ReadEnvironment(job.Result).
				Return(&models.Environment{}, nil)
		}

		c := getCLIContext(t, []string{"name", "2"}, flags)
		if err := command.update(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestEnvironmentLink(t *testing.T) {
	client, resolver, command, ctrl := initEnvCommandTest(t)
	defer ctrl.Finish()

	resolver.EXPECT().
		Resolve("environment", "name1").
		Return([]string{"id1"}, nil)

	resolver.EXPECT().
		Resolve("environment", "name2").
		Return([]string{"id2"}, nil)

	client.EXPECT().
		CreateLink("id1", "id2").
		Return(nil)

	c := getCLIContext(t, []string{"name1", "name2"}, nil)
	if err := command.link(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentUnlink(t *testing.T) {
	client, resolver, command, ctrl := initEnvCommandTest(t)
	defer ctrl.Finish()

	resolver.EXPECT().
		Resolve("environment", "name1").
		Return([]string{"id1"}, nil)

	resolver.EXPECT().
		Resolve("environment", "name2").
		Return([]string{"id2"}, nil)

	client.EXPECT().
		DeleteLink("id1", "id2").
		Return(nil)

	c := getCLIContext(t, []string{"name1", "name2"}, nil)
	if err := command.unlink(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentLink_duplicateEnvironmentID(t *testing.T) {
	_, resolver, command, ctrl := initEnvCommandTest(t)
	defer ctrl.Finish()

	resolver.EXPECT().
		Resolve("environment", gomock.Any()).
		Return([]string{"id1"}, nil).
		Times(2)

	c := getCLIContext(t, []string{"name1", "name2"}, nil)
	if err := command.link(c); err == nil {
		t.Fatal("error was nil!")
	}
}

func TestEnvironmentUnlink_duplicateEnvironmentID(t *testing.T) {
	_, resolver, command, ctrl := initEnvCommandTest(t)
	defer ctrl.Finish()

	resolver.EXPECT().
		Resolve("environment", gomock.Any()).
		Return([]string{"id1"}, nil).
		Times(2)

	c := getCLIContext(t, []string{"name1", "name2"}, nil)
	if err := command.unlink(c); err == nil {
		t.Fatal("error was nil!")
	}
}

func initEnvCommandTest(t *testing.T) (*mock_client.MockClient, *mock_resolver.MockResolver, *EnvironmentCommand, *gomock.Controller) {
	tc, ctrl := newTestCommand(t)
	return tc.Client, tc.Resolver, NewEnvironmentCommand(tc.Command()), ctrl
}
