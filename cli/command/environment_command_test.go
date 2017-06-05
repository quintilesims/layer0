package command

import (
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/urfave/cli"
	"testing"
)

func TestCreateEnvironment(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	file, close := tempFile(t, "user_data")
	defer close()

	tc.Client.EXPECT().
		CreateEnvironment("name", "m3.large", 2, []byte("user_data"), "linux", "ami").
		Return(&models.Environment{}, nil)

	flags := map[string]interface{}{
		"size":      "m3.large",
		"min-count": 2,
		"user-data": file.Name(),
		"os":        "linux",
		"ami":       "ami",
	}

	c := testutils.GetCLIContext(t, []string{"name"}, flags)
	if err := command.Create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateEnvironment_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Create(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestDeleteEnvironment(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		DeleteEnvironment("id").
		Return("jobid", nil)

	c := testutils.GetCLIContext(t, []string{"name"}, nil)
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEnvironmentWait(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		DeleteEnvironment("id").
		Return("jobid", nil)

	tc.Client.EXPECT().
		WaitForJob("jobid", testutils.TEST_TIMEOUT).
		Return(nil)

	c := testutils.GetCLIContext(t, []string{"name"}, map[string]interface{}{"wait": true})
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEnvironment_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Delete(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestGetEnvironment(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetEnvironment("id").
		Return(&models.Environment{}, nil)

	c := testutils.GetCLIContext(t, []string{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetEnvironment_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Get(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestListEnvironments(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	tc.Client.EXPECT().
		ListEnvironments().
		Return([]*models.EnvironmentSummary{}, nil)

	c := testutils.GetCLIContext(t, nil, nil)
	if err := command.List(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentSetMinCount(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		UpdateEnvironment("id", 2).
		Return(&models.Environment{}, nil)

	c := testutils.GetCLIContext(t, []string{"name", "2"}, nil)
	if err := command.SetMinCount(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentSetMinCount_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg":      testutils.GetCLIContext(t, nil, nil),
		"Missing COUNT arg":     testutils.GetCLIContext(t, []string{"name"}, nil),
		"Non-integer COUNT arg": testutils.GetCLIContext(t, []string{"name", "2w"}, nil),
	}

	for name, c := range contexts {
		if err := command.SetMinCount(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestEnvironmentLink(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "name1").
		Return([]string{"id1"}, nil)

	tc.Resolver.EXPECT().
		Resolve("environment", "name2").
		Return([]string{"id2"}, nil)

	tc.Client.EXPECT().
		CreateLink("id1", "id2").
		Return(nil)

	c := testutils.GetCLIContext(t, []string{"name1", "name2"}, nil)
	if err := command.Link(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentLink_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing SOURCE arg":      testutils.GetCLIContext(t, []string{}, nil),
		"Missing DESTINATION arg": testutils.GetCLIContext(t, []string{"name"}, nil),
	}

	for name, c := range contexts {
		if err := command.Link(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestEnvironmentLink_duplicateEnvironmentID(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve(gomock.Any(), gomock.Any()).
		Return([]string{"id1"}, nil).
		Times(2)

	c := testutils.GetCLIContext(t, []string{"name1", "name2"}, nil)
	if err := command.Link(c); err == nil {
		t.Fatal("error was nil!")
	}
}

func TestEnvironmentUnlink(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "name1").
		Return([]string{"id1"}, nil)

	tc.Resolver.EXPECT().
		Resolve("environment", "name2").
		Return([]string{"id2"}, nil)

	tc.Client.EXPECT().
		DeleteLink("id1", "id2").
		Return(nil)

	c := testutils.GetCLIContext(t, []string{"name1", "name2"}, nil)
	if err := command.Unlink(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentUnlink_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing SOURCE arg":      testutils.GetCLIContext(t, []string{}, nil),
		"Missing DESTINATION arg": testutils.GetCLIContext(t, []string{"name"}, nil),
	}

	for name, c := range contexts {
		if err := command.Unlink(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestEnvironmentUnlink_duplicateEnvironmentID(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", gomock.Any()).
		Return([]string{"id1"}, nil).
		Times(2)

	c := testutils.GetCLIContext(t, []string{"name1", "name2"}, nil)
	if err := command.Unlink(c); err == nil {
		t.Fatal("error was nil!")
	}
}
