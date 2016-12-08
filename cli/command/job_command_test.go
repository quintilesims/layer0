package command

import (
	"github.com/urfave/cli"
	"github.com/quintilesims/layer0/common/models"
	"testing"
)

func TestDeleteJob(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("job", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		DeleteJob("id").
		Return(nil)

	c := getCLIContext(t, Args{"name"}, nil)
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteJob_UserInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Delete(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestGetJob(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("job", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetJob("id").
		Return(&models.Job{}, nil)

	c := getCLIContext(t, Args{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetJob_UserInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Get(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestListJobs(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	tc.Client.EXPECT().
		ListJobs().
		Return([]*models.Job{}, nil)

	c := getCLIContext(t, nil, nil)
	if err := command.List(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetJobLogs(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("job", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetJob("id").
		Return(&models.Job{TaskID: "task-id"}, nil)

	tc.Client.EXPECT().
		GetTaskLogs("task-id", 100).
		Return([]*models.LogFile{}, nil)

	c := getCLIContext(t, Args{"name"}, Flags{"tail": 100})
	if err := command.Logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetJobLogs_UserInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Logs(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}
