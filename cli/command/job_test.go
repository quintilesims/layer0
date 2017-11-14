package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func TestJobDelete_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewJobCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.delete(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestJobRead_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewJobCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": NewContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.read(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestDeleteJob(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewJobCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("job", "job_name").
		Return([]string{"job_id"}, nil)

	base.Client.EXPECT().
		DeleteJob("job_id").
		Return(nil)

	c := NewContext(t, []string{"job_name"}, nil)
	if err := command.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadJob(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewJobCommand(base.Command())

	jobIDs := []string{"job_id1", "job_id2"}

	base.Resolver.EXPECT().
		Resolve("job", "job_*").
		Return(jobIDs, nil)

	for _, jobID := range jobIDs {
		base.Client.EXPECT().
			ReadJob(jobID).
			Return(&models.Job{}, nil)
	}

	c := NewContext(t, []string{"job_*"}, nil)
	if err := command.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestListJobs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewJobCommand(base.Command())

	base.Client.EXPECT().
		ListJobs().
		Return([]*models.Job{}, nil)

	c := NewContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}
