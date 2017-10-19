package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCreateTask(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()
		taskCommand := NewTaskCommand(base.Command())
		args := Args{"env_name", "dpl_name", "task_name"}
		flags := Flags{
			"copies": 1,
			"env":    []string{"container:key=val"},
		}
		overrides := []models.ContainerOverride{{
			ContainerName:        "container",
			EnvironmentOverrides: map[string]string{"key": "val"},
		}}

		req := models.CreateTaskRequest{
			TaskName:           "task_name",
			DeployID:           "dpl_id",
			EnvironmentID:      "env_id",
			ContainerOverrides: overrides,
		}

		base.Resolver.EXPECT().
			Resolve("deploy", "dpl_name").
			Return(Args{"dpl_id"}, nil)

		base.Resolver.EXPECT().
			Resolve("environment", "env_name").
			Return(Args{"env_id"}, nil)

		base.Client.EXPECT().
			CreateTask(req).
			Return("job_id", nil)

		if wait {
			base.Client.EXPECT().
				ReadJob("job_id").
				Return(&models.Job{
					Status: "Completed",
					Result: "task_id",
				}, nil)

			base.Client.EXPECT().
				ReadTask("task_id").
				Return(&models.Task{}, nil)
		}

		c := NewContext(t, args, flags, SetNoWait(!wait))
		if err := taskCommand.create(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCreateTask_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": NewContext(t, nil, nil),
		"Missing NAME arg":        NewContext(t, Args{"environment"}, nil),
		"Missing DEPLOY arg":      NewContext(t, Args{"environment", "name"}, nil),
	}

	for name, c := range contexts {
		if err := taskCommand.create(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestDeleteTask(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("task", "task_name").
		Return(Args{"task_id"}, nil)

	base.Client.EXPECT().
		DeleteTask("task_id").
		Return("job_id", nil)

	job := &models.Job{
		Status: "Completed",
		Result: "job_id",
	}

	base.Client.EXPECT().
		ReadJob("job_id").
		Return(job, nil)

	c := NewContext(t, Args{"task_name"}, nil)
	if err := taskCommand.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTask_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := taskCommand.delete(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestReadTask(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())
	result := []*models.TaskSummary{
		{TaskID: "task_id"},
	}

	base.Resolver.EXPECT().
		Resolve("task", "task_name").
		Return(Args{"task_id"}, nil)

	base.Client.EXPECT().
		ListTasks().
		Return(result, nil)

	base.Client.EXPECT().
		ReadTask("task_id").
		Return(&models.Task{}, nil)

	c := NewContext(t, Args{"task_name"}, nil)
	if err := taskCommand.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadTask_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())

	base.Client.EXPECT().
		ListTasks().
		Return([]*models.TaskSummary{}, nil)

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := taskCommand.read(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestReadTask_expiredTasks(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())
	result := []*models.TaskSummary{
		{TaskID: "task_id"},
	}

	base.Resolver.EXPECT().
		Resolve("task", "task_name").
		Return(Args{"task_id", "task_id2", "task_id3"}, nil)

	base.Client.EXPECT().
		ListTasks().
		Return(result, nil)

	base.Client.EXPECT().
		ReadTask("task_id").
		Return(&models.Task{}, nil)

	c := NewContext(t, Args{"task_name"}, nil)
	if err := taskCommand.read(c); err != nil {
		t.Fatal(err)
	}

}

func TestListTasks(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())

	base.Client.EXPECT().
		ListTasks().
		Return([]*models.TaskSummary{}, nil)

	c := NewContext(t, nil, nil)
	if err := taskCommand.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadTaskLogs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("task", "name").
		Return(Args{"id"}, nil)

	query := buildLogQueryHelper("id", "start", "end", 100)

	base.Client.EXPECT().
		ReadTaskLogs("id", query)

	flags := Flags{
		"tail":  100,
		"start": "start",
		"end":   "end",
	}

	c := NewContext(t, Args{"name"}, flags)
	if err := taskCommand.logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadTaskLogs_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := taskCommand.logs(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestParseOverrides(t *testing.T) {
	input := Args{
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

	assert.Equal(t, len(output), len(expected))
	assert.Equal(t, expected[0], output[0])
	assert.Equal(t, expected[1], output[1])
}
func TestParseOverridesErrors(t *testing.T) {
	cases := map[string]string{
		"Missing CONTAINER": ":key=val",
		"Missing KEY":       "container:=val",
		"Missing VAL":       "container:key=",
	}

	for name, input := range cases {
		if _, err := parseOverrides(Args{input}); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}
