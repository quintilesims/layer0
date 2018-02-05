package command

import (
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCreateTask(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	expected := []models.ContainerOverride{
		{
			ContainerName:        "c1",
			EnvironmentOverrides: map[string]string{"k1": "v1"},
		},
		{
			ContainerName:        "c2",
			EnvironmentOverrides: map[string]string{"k2": "v2"},
		},
	}

	validateOverride := func(req models.CreateTaskRequest) {
		assert.Equal(t, "tsk_name", req.TaskName)
		assert.Equal(t, "dpl_id", req.DeployID)
		assert.Equal(t, "env_id", req.EnvironmentID)

		assert.Len(t, req.ContainerOverrides, 2)
		for _, e := range expected {
			assert.Contains(t, req.ContainerOverrides, e)
		}
	}

	base.Client.EXPECT().
		CreateTask(gomock.Any()).
		Do(validateOverride).
		Return("tsk_id", nil)

	base.Client.EXPECT().
		ReadTask("tsk_id").
		Return(&models.Task{}, nil)

	flags := map[string]interface{}{
		"env": []string{
			"c1:k1=v1",
			"c2:k2=v2",
		},
	}

	c := testutils.NewTestContext(t, []string{"env_name", "tsk_name", "dpl_name"}, flags)
	if err := command.create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTaskInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing ENVIRONMENT arg": testutils.NewTestContext(t, nil, nil),
		"Missing NAME arg":        testutils.NewTestContext(t, []string{"env_name"}, nil),
		"Missing DEPLOY arg":      testutils.NewTestContext(t, []string{"env_name", "tsk_name"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.create(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("task", "tsk_name").
		Return([]string{"tsk_id"}, nil)

	base.Client.EXPECT().
		DeleteTask("tsk_id").
		Return(nil)

	c := testutils.NewTestContext(t, []string{"tsk_name"}, nil)
	if err := command.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTaskInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.delete(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestListTasks(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.Command())

	base.Client.EXPECT().
		ListTasks().
		Return([]models.TaskSummary{}, nil)

	c := testutils.NewTestContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadTask(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("task", "tsk_name").
		Return([]string{"tsk_id"}, nil)

	base.Client.EXPECT().
		ReadTask("tsk_id").
		Return(&models.Task{}, nil)

	c := testutils.NewTestContext(t, []string{"tsk_name"}, nil)
	if err := command.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadTaskInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.read(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestReadTaskLogs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.Command())

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

	flags := map[string]interface{}{
		"tail":  100,
		"start": "start",
		"end":   "end",
	}

	c := testutils.NewTestContext(t, []string{"tsk_name"}, flags)
	if err := command.logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadTaskLogsInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewTaskCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.logs(c); err == nil {
				t.Fatal("error was nil!")
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
		t.Run(name, func(t *testing.T) {
			if _, err := parseOverrides([]string{input}); err == nil {
				t.Fatalf("error was nil!")
			}
		})
	}
}
