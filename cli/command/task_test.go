package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCreateTask(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		taskCommand := NewTaskCommand(base.Command())

		args := []string{"env_name", "tsk_name", "dpl_name"}

		flags := map[string]interface{}{
			"copies": 1,
			"env":    []string{"container:key=val"},
		}

		overrides := []models.ContainerOverride{{
			ContainerName:        "container",
			EnvironmentOverrides: map[string]string{"key": "val"},
		}}

		req := models.CreateTaskRequest{
			TaskName:           "tsk_name",
			DeployID:           "dpl_id",
			EnvironmentID:      "env_id",
			ContainerOverrides: overrides,
		}

		base.Resolver.EXPECT().
			Resolve("deploy", "dpl_name").
			Return([]string{"dpl_id"}, nil)

		base.Resolver.EXPECT().
			Resolve("environment", "env_name").
			Return([]string{"env_id"}, nil)

		base.Client.EXPECT().
			CreateTask(req).
			Return("job_id", nil)

		if wait {
			job := &models.Job{
				Status: models.CompletedJobStatus,
				Result: "tsk_id",
			}

			base.Client.EXPECT().
				ReadJob("job_id").
				Return(job, nil)

			base.Client.EXPECT().
				ReadTask("tsk_id").
				Return(&models.Task{}, nil)
		}

		c := config.NewTestContext(t, args, flags, config.SetNoWait(!wait))
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
		"Missing ENVIRONMENT arg": config.NewTestContext(t, nil, nil),
		"Missing TASK_NAME arg":   config.NewTestContext(t, []string{"environment"}, nil),
		"Missing DEPLOY arg":      config.NewTestContext(t, []string{"environment", "name"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := taskCommand.create(c); err == nil {
				t.Fatalf("Error was nil!")
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		taskCommand := NewTaskCommand(base.Command())

		base.Resolver.EXPECT().
			Resolve("task", "task_name").
			Return([]string{"tsk_id"}, nil)

		base.Client.EXPECT().
			DeleteTask("tsk_id").
			Return("job_id", nil)

		if wait {
			job := &models.Job{
				Status: models.CompletedJobStatus,
				Result: "job_id",
			}

			base.Client.EXPECT().
				ReadJob("job_id").
				Return(job, nil)
		}

		c := config.NewTestContext(t, []string{"task_name"}, nil, config.SetNoWait(!wait))
		if err := taskCommand.delete(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDeleteTask_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	taskCommand := NewTaskCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing TASK_NAME arg": config.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := taskCommand.delete(c); err == nil {
				t.Fatalf("Error was nil!")
			}
		})
	}
}

func TestReadTask(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	taskCommand := NewTaskCommand(base.Command())

	result := []*models.TaskSummary{
		{TaskID: "tsk_id"},
	}

	base.Resolver.EXPECT().
		Resolve("task", "tsk_name").
		Return([]string{"tsk_id"}, nil)

	base.Client.EXPECT().
		ListTasks().
		Return(result, nil)

	base.Client.EXPECT().
		ReadTask("tsk_id").
		Return(&models.Task{}, nil)

	c := config.NewTestContext(t, []string{"tsk_name"}, nil)
	if err := taskCommand.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadTask_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	taskCommand := NewTaskCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing TASK_NAME arg": config.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := taskCommand.read(c); err == nil {
				t.Fatalf("Error was nil!")
			}
		})
	}
}

func TestReadTask_expiredTasks(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	taskCommand := NewTaskCommand(base.Command())

	result := []*models.TaskSummary{
		{TaskID: "tsk_id1"},
	}

	base.Resolver.EXPECT().
		Resolve("task", "tsk_name").
		Return([]string{"tsk_id1", "expired1", "expired2"}, nil)

	base.Client.EXPECT().
		ListTasks().
		Return(result, nil)

	base.Client.EXPECT().
		ReadTask("tsk_id1").
		Return(&models.Task{}, nil)

	c := config.NewTestContext(t, []string{"tsk_name"}, nil)
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

	c := config.NewTestContext(t, nil, nil)
	if err := taskCommand.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadTaskLogs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("task", "tsk_name").
		Return([]string{"tsk_id"}, nil)

	query := buildLogQueryHelper("start", "end", 100)

	base.Client.EXPECT().
		ReadTaskLogs("tsk_id", query).
		Return([]*models.LogFile{}, nil)

	flags := map[string]interface{}{
		"tail":  100,
		"start": "start",
		"end":   "end",
	}

	c := config.NewTestContext(t, []string{"tsk_name"}, flags)
	if err := taskCommand.logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadTaskLogs_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	taskCommand := NewTaskCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing TASK_NAME arg": config.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := taskCommand.logs(c); err == nil {
				t.Fatalf("Error was nil!")
			}
		})
	}
}

func TestParseOverrides(t *testing.T) {
	input := []string{
		"container1:key1=val1",
		"container1:key2=val2",
		"container2:k1=v1",
	}

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

	assert.Len(t, output, 2)
	for _, e := range expected {
		assert.Contains(t, output, e)
	}
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
