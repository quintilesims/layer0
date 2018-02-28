package command

import (
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	validateRequest := func(req models.CreateTaskRequest) {
		assert.Equal(t, "tsk_name", req.TaskName)
		assert.Equal(t, "env_id", req.EnvironmentID)
		assert.Equal(t, "dpl_id", req.DeployID)

		overrides := []models.ContainerOverride{
			{ContainerName: "c1", EnvironmentOverrides: map[string]string{"k1": "v1"}},
			{ContainerName: "c2", EnvironmentOverrides: map[string]string{"k2": "v2"}},
		}

		assert.Len(t, req.ContainerOverrides, len(overrides))
		assert.Contains(t, req.ContainerOverrides, overrides[0])
		assert.Contains(t, req.ContainerOverrides, overrides[1])
	}

	base.Client.EXPECT().
		CreateTask(gomock.Any()).
		Do(validateRequest).
		Return("tsk_id", nil)

	base.Client.EXPECT().
		ReadTask("tsk_id").
		Return(&models.Task{}, nil)

	input := "l0 task create "
	input += "--env c1:k1=v1 "
	input += "--env c2:k2=v2 "
	input += "env_name tsk_name dpl_name"

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTask_stateful(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	overrides := []models.ContainerOverride{
		{ContainerName: "c1", EnvironmentOverrides: map[string]string{"k1": "v1"}},
		{ContainerName: "c2", EnvironmentOverrides: map[string]string{"k2": "v2"}},
	}

	validateRequest := func(req models.CreateTaskRequest) {
		assert.Equal(t, "tsk_name", req.TaskName)
		assert.Equal(t, "env_id", req.EnvironmentID)
		assert.Equal(t, "dpl_id", req.DeployID)

		assert.Len(t, req.ContainerOverrides, len(overrides))
		assert.Contains(t, req.ContainerOverrides, overrides[0])
		assert.Contains(t, req.ContainerOverrides, overrides[1])
	}

	base.Client.EXPECT().
		CreateTask(gomock.Any()).
		Do(validateRequest).
		Return("tsk_id", nil)

	base.Client.EXPECT().
		ReadTask("tsk_id").
		Return(&models.Task{}, nil)

	input := "l0 task create "
	input += "--stateful "
	input += "--env c1:k1=v1 "
	input += "--env c2:k2=v2 "
	input += "env_name tsk_name dpl_name"

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTaskInputErrors(t *testing.T) {
	testInputErrors(t, NewTaskCommand(nil).Command(), map[string]string{
		"Missing ENVIRONMENT arg": "l0 task create",
		"Missing NAME arg":        "l0 task create env",
		"Missing DEPLOY arg":      "l0 task create env name",
	})
}

func TestDeleteTask(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("task", "tsk_name").
		Return([]string{"tsk_id"}, nil)

	base.Client.EXPECT().
		DeleteTask("tsk_id").
		Return(nil)

	input := "l0 task delete tsk_name"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTaskInputErrors(t *testing.T) {
	testInputErrors(t, NewTaskCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 task delete",
	})
}

func TestListTasks(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.CommandBase()).Command()

	base.Client.EXPECT().
		ListTasks().
		Return([]models.TaskSummary{}, nil)

	input := "l0 task list"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestReadTask(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("task", "tsk_name").
		Return([]string{"tsk_id"}, nil)

	base.Client.EXPECT().
		ReadTask("tsk_id").
		Return(&models.Task{}, nil)

	input := "l0 task get tsk_name"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestReadTaskInputErrors(t *testing.T) {
	testInputErrors(t, NewTaskCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 task get",
	})
}

func TestReadTaskLogs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("task", "tsk_name").
		Return([]string{"tsk_id"}, nil)

	query := url.Values{
		"tail":  []string{"100"},
		"start": []string{"start"},
		"end":   []string{"end"},
	}

	base.Client.EXPECT().
		ReadTaskLogs("tsk_id", query).
		Return([]models.LogFile{}, nil)

	input := "l0 task logs "
	input += "--tail 100 "
	input += "--start start "
	input += "--end end "
	input += "tsk_name"

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestReadTaskLogsInputErrors(t *testing.T) {
	testInputErrors(t, NewTaskCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 task logs",
	})
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
		t.Run(name, func(t *testing.T) {
			if _, err := parseOverrides([]string{input}); err == nil {
				t.Fatalf("error was nil!")
			}
		})
	}
}
