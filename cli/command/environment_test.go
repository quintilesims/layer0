package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func TestEnvironmentCreate_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg":  config.NewTestContext(t, nil, nil),
		"Negative MinScale": config.NewTestContext(t, []string{"env_name"}, map[string]interface{}{"min-scale": "-1"}),
		"Negative MaxScale": config.NewTestContext(t, []string{"env_name"}, map[string]interface{}{"max-scale": "-1"}),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.create(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestEnvironmentDelete_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": config.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.delete(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestEnvironmentRead_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": config.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.read(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestEnvironmentSetScale_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": config.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.setScale(c); err == nil {
				t.Fatalf("%s: error was nil!", name)
			}
		})
	}
}

func TestEnvironmentLink_duplicateEnvironmentID(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", gomock.Any()).
		Return([]string{"env_id1"}, nil).
		Times(2)

	c := config.NewTestContext(t, []string{"env_name1", "env_name1"}, nil)
	if err := command.link(c); err == nil {
		t.Fatal("error was nil!")
	}
}

func TestEnvironmentUnlink_duplicateEnvironmentID(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", gomock.Any()).
		Return([]string{"env_id1"}, nil).
		Times(2)

	c := config.NewTestContext(t, []string{"env_name1", "env_name1"}, nil)
	if err := command.unlink(c); err == nil {
		t.Fatal("error was nil!")
	}
}

func TestCreateEnvironment(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		userData := "user_data"
		file, close := createTempFile(t, userData)
		defer close()

		req := models.CreateEnvironmentRequest{
			EnvironmentName:  "env_name",
			InstanceType:     "t2.small",
			MinScale:         2,
			MaxScale:         5,
			UserDataTemplate: []byte(userData),
			OperatingSystem:  "linux",
			AMIID:            "ami",
		}

		environment := &models.Environment{}
		job := &models.Job{
			JobID:  "job_id",
			Status: models.CompletedJobStatus,
			Result: "entity_id",
		}

		base.Client.EXPECT().
			CreateEnvironment(req).
			Return(job.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)

			base.Client.EXPECT().
				ReadEnvironment(job.Result).
				Return(environment, nil)
		}

		flags := map[string]interface{}{
			"type":      req.InstanceType,
			"min-scale": req.MinScale,
			"max-scale": req.MaxScale,
			"user-data": file.Name(),
			"os":        req.OperatingSystem,
			"ami":       req.AMIID,
		}

		c := config.NewTestContext(t, []string{"env_name"}, flags, config.SetNoWait(!wait))
		if err := command.create(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDeleteEnvironment(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		job := &models.Job{
			JobID:  "job_id",
			Status: models.CompletedJobStatus,
			Result: "entity_id",
		}

		base.Resolver.EXPECT().
			Resolve("environment", "env_name").
			Return([]string{"env_id"}, nil)

		base.Client.EXPECT().
			DeleteEnvironment("env_id").
			Return(job.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)
		}

		c := config.NewTestContext(t, []string{"env_name"}, nil, config.SetNoWait(!wait))
		if err := command.delete(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetEnvironment(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", "env_name*").
		Return([]string{"env_id1", "env_id2"}, nil)

	base.Client.EXPECT().
		ReadEnvironment("env_id1").
		Return(&models.Environment{}, nil)

	base.Client.EXPECT().
		ReadEnvironment("env_id2").
		Return(&models.Environment{}, nil)

	c := config.NewTestContext(t, []string{"env_name*"}, nil)
	if err := command.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestListEnvironments(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	command := NewEnvironmentCommand(base.Command())

	base.Client.EXPECT().
		ListEnvironments().
		Return([]*models.EnvironmentSummary{}, nil)

	c := config.NewTestContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentSetScale(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		base.Resolver.EXPECT().
			Resolve("environment", "env_name").
			Return([]string{"env_id"}, nil)

		minScale := 2
		maxScale := 5
		req := models.UpdateEnvironmentRequest{
			MinScale: &minScale,
			MaxScale: &maxScale,
		}

		job := &models.Job{
			JobID:  "job_id",
			Status: models.CompletedJobStatus,
			Result: "entity_id",
		}

		base.Client.EXPECT().
			UpdateEnvironment("env_id", req).
			Return(job.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job.JobID).
				Return(job, nil)

			base.Client.EXPECT().
				ReadEnvironment(job.Result).
				Return(&models.Environment{}, nil)
		}

		flags := map[string]interface{}{
			"min-scale": 2,
			"max-scale": 5,
		}

		c := config.NewTestContext(t, []string{"env_name"}, flags, config.SetNoWait(!wait))
		if err := command.setScale(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestEnvironmentLinkBiDirectional(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		links1 := []string{"env_id2"}
		links2 := []string{"env_id3", "env_id1"}
		env1 := &models.Environment{
			EnvironmentID: "env_id1",
			Links:         []string{},
		}

		env2 := &models.Environment{
			EnvironmentID: "env_id2",
			Links:         []string{"env_id3"},
		}

		base.Resolver.EXPECT().
			Resolve("environment", "env_name1").
			Return([]string{"env_id1"}, nil)

		base.Resolver.EXPECT().
			Resolve("environment", "env_name2").
			Return([]string{"env_id2"}, nil)

		base.Client.EXPECT().
			ReadEnvironment("env_id1").
			Return(env1, nil)

		base.Client.EXPECT().
			ReadEnvironment("env_id2").
			Return(env2, nil)

		job1 := &models.Job{
			JobID:  "job_id1",
			Status: models.CompletedJobStatus,
			Result: "env_id1",
		}

		req1 := models.UpdateEnvironmentRequest{Links: &links1}

		base.Client.EXPECT().
			UpdateEnvironment("env_id1", req1).
			Return(job1.JobID, nil)

		job2 := &models.Job{
			JobID:  "job_id2",
			Status: models.CompletedJobStatus,
			Result: "env_id2",
		}

		req2 := models.UpdateEnvironmentRequest{Links: &links2}

		base.Client.EXPECT().
			UpdateEnvironment("env_id2", req2).
			Return(job2.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job1.JobID).
				Return(job1, nil)

			base.Client.EXPECT().
				ReadJob(job2.JobID).
				Return(job2, nil)
		}

		f := map[string]interface{}{
			"bi-directional": true,
		}

		c := config.NewTestContext(t, []string{"env_name1", "env_name2"}, f, config.SetNoWait(!wait))
		if err := command.link(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestEnvironmentLinkUniDirectional(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		links1 := []string{"env_id3", "env_id2"}
		env1 := &models.Environment{
			EnvironmentID: "env_id1",
			Links:         []string{"env_id3"},
		}

		base.Resolver.EXPECT().
			Resolve("environment", "env_name1").
			Return([]string{"env_id1"}, nil)

		base.Resolver.EXPECT().
			Resolve("environment", "env_name2").
			Return([]string{"env_id2"}, nil)

		base.Client.EXPECT().
			ReadEnvironment("env_id1").
			Return(env1, nil)

		job1 := &models.Job{
			JobID:  "job_id1",
			Status: models.CompletedJobStatus,
			Result: "env_id1",
		}

		req1 := models.UpdateEnvironmentRequest{Links: &links1}

		base.Client.EXPECT().
			UpdateEnvironment("env_id1", req1).
			Return(job1.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job1.JobID).
				Return(job1, nil)
		}

		c := config.NewTestContext(t, []string{"env_name1", "env_name2"}, nil, config.SetNoWait(!wait))
		if err := command.link(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestEnvironmentUnlinkBiDirectional(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())

		links1 := []string{"env_id3"}
		links2 := []string{"env_id3"}
		env1 := &models.Environment{
			EnvironmentID: "env_id1",
			Links:         []string{"env_id2", "env_id3"},
		}
		env2 := &models.Environment{
			EnvironmentID: "env_id2",
			Links:         []string{"env_id1", "env_id3"},
		}

		base.Resolver.EXPECT().
			Resolve("environment", "env_name1").
			Return([]string{"env_id1"}, nil)

		base.Resolver.EXPECT().
			Resolve("environment", "env_name2").
			Return([]string{"env_id2"}, nil)

		base.Client.EXPECT().
			ReadEnvironment("env_id1").
			Return(env1, nil)

		base.Client.EXPECT().
			ReadEnvironment("env_id2").
			Return(env2, nil)

		job1 := &models.Job{
			JobID:  "job_id1",
			Status: models.CompletedJobStatus,
			Result: "env_id1",
		}

		req1 := models.UpdateEnvironmentRequest{Links: &links1}

		base.Client.EXPECT().
			UpdateEnvironment("env_id1", req1).
			Return(job1.JobID, nil)

		job2 := &models.Job{
			JobID:  "job_id2",
			Status: models.CompletedJobStatus,
			Result: "env_id2",
		}

		req2 := models.UpdateEnvironmentRequest{Links: &links2}

		base.Client.EXPECT().
			UpdateEnvironment("env_id2", req2).
			Return(job2.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job1.JobID).
				Return(job1, nil)

			base.Client.EXPECT().
				ReadJob(job2.JobID).
				Return(job2, nil)
		}

		f := map[string]interface{}{
			"bi-directional": true,
		}

		c := config.NewTestContext(t, []string{"env_name1", "env_name2"}, f, config.SetNoWait(!wait))
		if err := command.unlink(c); err != nil {
			t.Fatal(err)
		}
	})
}

func TestEnvironmentUnlinkUnidirectional(t *testing.T) {
	testWaitHelper(t, func(t *testing.T, wait bool) {
		base, ctrl := newTestCommand(t)
		defer ctrl.Finish()

		command := NewEnvironmentCommand(base.Command())
		links1 := []string{}
		env1 := &models.Environment{
			EnvironmentID: "env_id1",
			Links:         []string{"env_id2"},
		}

		base.Resolver.EXPECT().
			Resolve("environment", "env_name1").
			Return([]string{"env_id1"}, nil)

		base.Resolver.EXPECT().
			Resolve("environment", "env_name2").
			Return([]string{"env_id2"}, nil)

		base.Client.EXPECT().
			ReadEnvironment("env_id1").
			Return(env1, nil)

		job1 := &models.Job{
			JobID:  "job_id1",
			Status: models.CompletedJobStatus,
			Result: "env_id1",
		}

		req1 := models.UpdateEnvironmentRequest{Links: &links1}

		base.Client.EXPECT().
			UpdateEnvironment("env_id1", req1).
			Return(job1.JobID, nil)

		if wait {
			base.Client.EXPECT().
				ReadJob(job1.JobID).
				Return(job1, nil)
		}

		c := config.NewTestContext(t, []string{"env_name1", "env_name2"}, nil, config.SetNoWait(!wait))
		if err := command.unlink(c); err != nil {
			t.Fatal(err)
		}
	})
}
