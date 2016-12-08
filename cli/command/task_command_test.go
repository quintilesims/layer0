package command

import (
	"github.com/urfave/cli"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestParseOverrides(t *testing.T) {
	input := []string{
		"container1:key1=val1",
		"container1:key2=val2",
		"container2:k1=v1"}

	expected := []models.ContainerOverride{
		{
			ContainerName:        "container1",
			EnvironmentOverrides: map[string]string{"key1": "val1", "key2": "val2"},
		},
		{
			ContainerName:        "container2",
			EnvironmentOverrides: map[string]string{"k1": "v1"},
		},
	}

	output, err := parseOverrides(input)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(output), len(expected))
	testutils.AssertInSlice(t, expected[0], output)
	testutils.AssertInSlice(t, expected[1], output)
}

func TestParseOverridesErrors(t *testing.T) {
	cases := map[string]string{
		"Missing CONTAINER": ":key=val",
		"Missing KEY":       "container:=val",
		"Missing VAL":       "container:key=",
	}

	for name, input := range cases {
		if _, err := parseOverrides([]string{input}); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestCreateTask(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "environment").
		Return([]string{"environmentID"}, nil)

	tc.Resolver.EXPECT().
		Resolve("deploy", "deploy").
		Return([]string{"deployID"}, nil)

	overrides := []models.ContainerOverride{{
		ContainerName:        "container",
		EnvironmentOverrides: map[string]string{"key": "val"},
	}}

	tc.Client.EXPECT().
		CreateTask("name", "environmentID", "deployID", 2, overrides).
		Return(&models.Task{}, nil)

	flags := Flags{
		"copies": 2,
		"env":    []string{"container:key=val"},
	}

	c := getCLIContext(t, Args{"environment", "name", "deploy"}, flags)
	if err := command.Create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTask_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": getCLIContext(t, nil, nil),
		"Missing NAME arg":        getCLIContext(t, Args{"environment"}, nil),
		"Missing DEPLOY arg":      getCLIContext(t, Args{"environment", "name"}, nil),
	}

	for name, c := range contexts {
		if err := command.Create(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestDeleteTask(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("task", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		DeleteTask("id").
		Return(nil)

	c := getCLIContext(t, Args{"name"}, nil)
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTask_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Delete(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestGetTask(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("task", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetTask("id").
		Return(&models.Task{}, nil)

	c := getCLIContext(t, Args{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetTask_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Get(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestListTasks(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	tc.Client.EXPECT().
		ListTasks().
		Return([]*models.Task{}, nil)

	c := getCLIContext(t, nil, nil)
	if err := command.List(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetTaskLogs(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("task", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetTaskLogs("id", 100)

	c := getCLIContext(t, Args{"name"}, Flags{"tail": 100})
	if err := command.Logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetTaskLogs_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Logs(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}
