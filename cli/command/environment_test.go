package command

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/cli/printer"
	"github.com/quintilesims/layer0/cli/resolver/mock_resolver"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func TestCreateEnvironment_userInputErrors(t *testing.T) {
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
			c := GetCLIContext(t, tc.args, nil)

			if err := tc.command(c); err == nil {
				t.Fatalf("%s: error was nil!", tc.name)
			}
		})
	}
}

func TestCreateEnvironment(t *testing.T) {
	testCases := []struct {
		name  string
		flags map[string]interface{}
		wait  bool
	}{
		{
			name: "Wait",
			wait: true,
		},
		{
			name:  "NoWait",
			flags: map[string]interface{}{"no-wait": true},
			wait:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			CreateEnvironmentNoWait(t, tc.flags, tc.wait)
		})
	}
}

func CreateEnvironmentNoWait(t *testing.T, otherFlags map[string]interface{}, wait bool) {
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

	c := GetCLIContext(t, []string{"name"}, flags)
	if err := command.create(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEnvironment(t *testing.T) {
	testCases := []struct {
		name  string
		flags map[string]interface{}
		wait  bool
	}{
		{
			name: "Wait",
			wait: true,
		},
		{
			name:  "NoWait",
			flags: map[string]interface{}{"no-wait": true},
			wait:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

			if tc.wait {
				client.EXPECT().
					ReadJob(job.JobID).
					Return(job, nil)
			}

			c := GetCLIContext(t, []string{"name"}, tc.flags)
			if err := command.delete(c); err != nil {
				t.Fatal(err)
			}
		})
	}
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

	c := GetCLIContext(t, []string{"name"}, nil)
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

	c := GetCLIContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentSetMinCount(t *testing.T) {
	testCases := []struct {
		name  string
		flags map[string]interface{}
		wait  bool
	}{
		{
			name: "Wait",
			wait: true,
		},
		{
			name:  "NoWait",
			flags: map[string]interface{}{"no-wait": true},
			wait:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

			if tc.wait {
				client.EXPECT().
					ReadJob(job.JobID).
					Return(job, nil)

				client.EXPECT().
					ReadEnvironment(job.Result).
					Return(&models.Environment{}, nil)
			}

			c := GetCLIContext(t, []string{"name", "2"}, tc.flags)
			if err := command.update(c); err != nil {
				t.Fatal(err)
			}
		})
	}
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

	c := GetCLIContext(t, []string{"name1", "name2"}, nil)
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

	c := GetCLIContext(t, []string{"name1", "name2"}, nil)
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

	c := GetCLIContext(t, []string{"name1", "name2"}, nil)
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

	c := GetCLIContext(t, []string{"name1", "name2"}, nil)
	if err := command.unlink(c); err == nil {
		t.Fatal("error was nil!")
	}
}

func initEnvCommandTest(t *testing.T) (*mock_client.MockClient, *mock_resolver.MockResolver, *EnvironmentCommand, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	tc := &TestCommandBase{
		Client:   mock_client.NewMockClient(ctrl),
		Printer:  &printer.TestPrinter{},
		Resolver: mock_resolver.NewMockResolver(ctrl),
	}

	envCmd := NewEnvironmentCommand(tc.Command())

	return tc.Client, tc.Resolver, envCmd, ctrl
}

func tempFile(t *testing.T, content string) (*os.File, func()) {
	file, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	return file, func() { os.Remove(file.Name()) }
}
