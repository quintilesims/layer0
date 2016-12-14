package command

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
	"testing"
)

func TestDelete(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("job", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		Delete("id").
		Return(nil)

	c := getCLIContext(t, Args{"name"}, nil)
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDelete_UserInputErrors(t *testing.T) {
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

func TestSelectByID(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("job", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		SelectByID("id").
		Return(&models.Job{}, nil)

	c := getCLIContext(t, Args{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestSelectByID_UserInputErrors(t *testing.T) {
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

func TestSelectAll(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	tc.Client.EXPECT().
		SelectAll().
		Return([]*models.Job{}, nil)

	c := getCLIContext(t, nil, nil)
	if err := command.List(c); err != nil {
		t.Fatal(err)
	}
}

func TestSelectByIDLogs(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewJobCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("job", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		SelectByID("id").
		Return(&models.Job{TaskID: "task-id"}, nil)

	tc.Client.EXPECT().
		GetTaskLogs("task-id", 100).
		Return([]*models.LogFile{}, nil)

	c := getCLIContext(t, Args{"name"}, Flags{"tail": 100})
	if err := command.Logs(c); err != nil {
		t.Fatal(err)
	}
}

func TestSelectByIDLogs_UserInputErrors(t *testing.T) {
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
