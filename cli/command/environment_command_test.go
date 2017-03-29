package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func TestCreateEnvironment(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	file, close := tempFile(t, "user_data")
	defer close()

	tc.Client.EXPECT().
		CreateEnvironment("name", "m3.large", 2, []byte("user_data"), "linux").
		Return(&models.Environment{}, nil)

	flags := Flags{
		"size":      "m3.large",
		"min-count": 2,
		"user-data": file.Name(),
		"os":        "linux",
	}

	c := getCLIContext(t, Args{"name"}, flags)
	if err := command.Create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateEnvironment_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
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

	c := getCLIContext(t, Args{"name"}, nil)
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
		WaitForJob("jobid", TEST_TIMEOUT).
		Return(nil)

	c := getCLIContext(t, Args{"name"}, Flags{"wait": true})
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEnvironment_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
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

	c := getCLIContext(t, Args{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetEnvironment_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
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

	c := getCLIContext(t, nil, nil)
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

	c := getCLIContext(t, Args{"name", "2"}, nil)
	if err := command.SetMinCount(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentSetMinCount_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg":      getCLIContext(t, nil, nil),
		"Missing COUNT arg":     getCLIContext(t, Args{"name"}, nil),
		"Non-integer COUNT arg": getCLIContext(t, Args{"name", "2w"}, nil),
	}

	for name, c := range contexts {
		if err := command.SetMinCount(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}
