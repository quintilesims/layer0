package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/urfave/cli"
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
		Return("jobid", nil)

	flags := map[string]interface{}{
		"copies": 2,
		"env":    []string{"container:key=val"},
	}

	c := testutils.GetCLIContext(t, []string{"environment", "name", "deploy"}, flags)
	if err := command.Create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTaskWait(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("environment", "environment").
		Return([]string{"environmentID"}, nil)

	tc.Resolver.EXPECT().
		Resolve("deploy", "deploy").
		Return([]string{"deployID"}, nil)

	tc.Client.EXPECT().
		CreateTask("name", "environmentID", "deployID", 0, []models.ContainerOverride{}).
		Return("jobid", nil)

	tc.Client.EXPECT().
		WaitForJob("jobid", gomock.Any()).
		Return(nil)

	jobMeta := map[string]string{"task_0": "tid0", "task_1": "tid1"}
	tc.Client.EXPECT().
		GetJob("jobid").
		Return(&models.Job{Meta: jobMeta}, nil)

	tc.Client.EXPECT().
		GetTask("tid0").
		Return(&models.Task{}, nil)

	tc.Client.EXPECT().
		GetTask("tid1").
		Return(&models.Task{}, nil)

	flags := map[string]interface{}{
		"wait": true,
	}

	c := testutils.GetCLIContext(t, []string{"environment", "name", "deploy"}, flags)
	if err := command.Create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTask_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": testutils.GetCLIContext(t, nil, nil),
		"Missing NAME arg":        testutils.GetCLIContext(t, []string{"environment"}, nil),
		"Missing DEPLOY arg":      testutils.GetCLIContext(t, []string{"environment", "name"}, nil),
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

	c := testutils.GetCLIContext(t, []string{"name"}, nil)
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTask_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
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
	result := []*models.TaskSummary{
		{TaskID: "id"},
	}

	tc.Resolver.EXPECT().
		Resolve("task", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		ListTasks().
		Return(result, nil)

	tc.Client.EXPECT().
		GetTask("id").
		Return(&models.Task{}, nil)

	c := testutils.GetCLIContext(t, []string{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetTask_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
	}

	tc.Client.EXPECT().
		ListTasks().
		Return(nil, nil)

	for name, c := range contexts {
		if err := command.Get(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestGetTask_expiredTasks(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())
	result := []*models.TaskSummary{
		{TaskID: "id3"},
	}

	tc.Resolver.EXPECT().
		Resolve("task", "name").
		Return([]string{"id3", "id4", "id5"}, nil)

	tc.Client.EXPECT().
		ListTasks().
		Return(result, nil)

	//only task 'id3' should result in a GetTask call
	tc.Client.EXPECT().
		GetTask("id3").
		Return(&models.Task{}, nil)

	c := testutils.GetCLIContext(t, []string{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestListTasks(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	tc.Client.EXPECT().
		ListTasks().
		Return([]*models.TaskSummary{}, nil)

	c := testutils.GetCLIContext(t, nil, nil)
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
		GetTaskLogs("id", "start", "end", 100)

	flags := map[string]interface{}{
		"tail":  100,
		"start": "start",
		"end":   "end",
	}

	c := testutils.GetCLIContext(t, []string{"name"}, flags)
	if err := command.Logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetTaskLogs_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.GetCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Logs(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestFilterTaskSummaries(t *testing.T) {
	input := []*models.TaskSummary{
		{TaskName: "a", TaskID: "a1"},
		{TaskName: "b", TaskID: "b1"},
		{TaskID: "nameless1"},
		{TaskID: "nameless2"},
	}

	output := filterTaskSummaries(input)

	testutils.AssertEqual(t, len(output), 2)
	// only 'a' and 'b' tasks
	testutils.AssertInSlice(t, input[0], output)
	testutils.AssertInSlice(t, input[1], output)
}
