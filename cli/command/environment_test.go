package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/urfave/cli"
)

func TestCreateEnvironment(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	file, delete := testutils.TempFile(t, "user_data")
	defer delete()

	req := models.CreateEnvironmentRequest{
		EnvironmentName:  "env_name",
		EnvironmentType:  "static",
		InstanceType:     "t2.small",
		Scale:            2,
		UserDataTemplate: []byte("user_data"),
		OperatingSystem:  "linux",
		AMIID:            "ami",
	}

	base.Client.EXPECT().
		CreateEnvironment(req).
		Return("env_id", nil)

	base.Client.EXPECT().
		ReadEnvironment("env_id").
		Return(&models.Environment{}, nil)

	flags := map[string]interface{}{
		"type":      req.InstanceType,
		"scale":     req.Scale,
		"user-data": file.Name(),
		"os":        req.OperatingSystem,
		"ami":       req.AMIID,
	}

	c := testutils.NewTestContext(t, []string{"env_name"}, flags)
	if err := command.create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateEnvironmentInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
		"Negative Scale": testutils.NewTestContext(t,
			[]string{"env_name"},
			map[string]interface{}{"scale": "-1"}),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.create(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestDeleteEnvironment(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Client.EXPECT().
		DeleteEnvironment("env_id").
		Return(nil)

	c := testutils.NewTestContext(t, []string{"env_name"}, nil)
	if err := command.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEnvironmentInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

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

func TestListEnvironments(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	base.Client.EXPECT().
		ListEnvironments().
		Return([]models.EnvironmentSummary{}, nil)

	c := testutils.NewTestContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadEnvironment(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Client.EXPECT().
		ReadEnvironment("env_id").
		Return(&models.Environment{}, nil)

	c := testutils.NewTestContext(t, []string{"env_name"}, nil)
	if err := command.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadEnvironmentInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

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

func TestLinkEnvironments(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	environments := map[string]string{
		"env_name1": "env_id1",
		"env_name2": "env_id2",
	}

	for name, id := range environments {
		base.Resolver.EXPECT().
			Resolve("environment", name).
			Return([]string{id}, nil)

		environment := &models.Environment{
			EnvironmentID: id,
			Links:         []string{},
		}

		base.Client.EXPECT().
			ReadEnvironment(id).
			Return(environment, nil)
	}

	links := []string{"env_id2"}
	req := models.UpdateEnvironmentRequest{
		Links: &links,
	}

	base.Client.EXPECT().
		UpdateEnvironment("env_id1", req).
		Return(nil)

	c := testutils.NewTestContext(t, []string{"env_name1", "env_name2"}, nil)
	if err := command.link(c); err != nil {
		t.Fatal(err)
	}
}

func TestLinkEnvironmentsBidirectional(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	environments := map[string]string{
		"env_name1": "env_id1",
		"env_name2": "env_id2",
	}

	for name, id := range environments {
		base.Resolver.EXPECT().
			Resolve("environment", name).
			Return([]string{id}, nil)

		environment := &models.Environment{
			EnvironmentID: id,
			Links:         []string{},
		}

		base.Client.EXPECT().
			ReadEnvironment(id).
			Return(environment, nil)

		var links []string
		if id == "env_id1" {
			links = []string{"env_id2"}
		} else {
			links = []string{"env_id1"}
		}

		req := models.UpdateEnvironmentRequest{
			Links: &links,
		}

		base.Client.EXPECT().
			UpdateEnvironment(id, req).
			Return(nil)
	}

	flags := map[string]interface{}{
		"bi-directional": true,
	}

	c := testutils.NewTestContext(t, []string{"env_name1", "env_name2"}, flags)
	if err := command.link(c); err != nil {
		t.Fatal(err)
	}
}

func TestLinkEnvironmentsInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing SOURCE arg": testutils.NewTestContext(t, nil, nil),
		"Missing DEST arg":   testutils.NewTestContext(t, []string{"env_name1"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.link(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestUnlinkEnvironments(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	environments := map[string]string{
		"env_name1": "env_id1",
		"env_name2": "env_id2",
	}

	for name, id := range environments {
		base.Resolver.EXPECT().
			Resolve("environment", name).
			Return([]string{id}, nil)

		var links []string
		if id == "env_id1" {
			links = []string{"env_id2"}
		} else {
			links = []string{"env_id1"}
		}

		environment := &models.Environment{
			EnvironmentID: id,
			Links:         links,
		}

		base.Client.EXPECT().
			ReadEnvironment(id).
			Return(environment, nil)
	}

	links := []string{}
	req := models.UpdateEnvironmentRequest{
		Links: &links,
	}

	base.Client.EXPECT().
		UpdateEnvironment("env_id1", req).
		Return(nil)

	c := testutils.NewTestContext(t, []string{"env_name1", "env_name2"}, nil)
	if err := command.unlink(c); err != nil {
		t.Fatal(err)
	}
}

func TestUnlinkEnvironmentsBidirectional(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	environments := map[string]string{
		"env_name1": "env_id1",
		"env_name2": "env_id2",
	}

	for name, id := range environments {
		base.Resolver.EXPECT().
			Resolve("environment", name).
			Return([]string{id}, nil)

		var links []string
		if id == "env_id1" {
			links = []string{"env_id2"}
		} else {
			links = []string{"env_id1"}
		}

		environment := &models.Environment{
			EnvironmentID: id,
			Links:         links,
		}

		base.Client.EXPECT().
			ReadEnvironment(id).
			Return(environment, nil)

		updatedLinks := []string{}
		req := models.UpdateEnvironmentRequest{
			Links: &updatedLinks,
		}

		base.Client.EXPECT().
			UpdateEnvironment(id, req).
			Return(nil)
	}

	flags := map[string]interface{}{
		"bi-directional": true,
	}

	c := testutils.NewTestContext(t, []string{"env_name1", "env_name2"}, flags)
	if err := command.unlink(c); err != nil {
		t.Fatal(err)
	}
}

func TestUnlinkEnvironmentsInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing SOURCE arg": testutils.NewTestContext(t, nil, nil),
		"Missing DEST arg":   testutils.NewTestContext(t, []string{"env_name1"}, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.unlink(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}
