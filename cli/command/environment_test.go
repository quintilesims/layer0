package command

import (
	"fmt"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestCreateEnvironment(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	file, delete := testutils.TempFile(t, "user_data")
	defer delete()

	req := models.CreateEnvironmentRequest{
		EnvironmentName:  "env_name",
		InstanceType:     "t2.small",
		MinScale:         2,
		MaxScale:         5,
		UserDataTemplate: []byte("user_data"),
		OperatingSystem:  "linux",
		AMIID:            "ami123",
	}

	base.Client.EXPECT().
		CreateEnvironment(req).
		Return("env_id", nil)

	base.Client.EXPECT().
		ReadEnvironment("env_id").
		Return(&models.Environment{}, nil)

	input := "l0 environment create "
	input += "--type t2.small "
	input += "--min-scale 2 "
	input += "--max-scale 5 "
	input += "--os linux "
	input += "--ami ami123 "
	input += fmt.Sprintf("--user-data %s ", file.Name())
	input += "env_name "

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestCreateEnvironmentInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	cases := map[string]string{
		"Missing NAME arg": "l0 environment create",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := testutils.RunApp(command, input); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestDeleteEnvironment(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Client.EXPECT().
		DeleteEnvironment("env_id").
		Return(nil)

	if err := testutils.RunApp(command, "l0 environment delete env_name"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEnvironmentInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	cases := map[string]string{
		"Missing NAME arg": "l0 environment delete",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := testutils.RunApp(command, input); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestListEnvironments(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	base.Client.EXPECT().
		ListEnvironments().
		Return([]models.EnvironmentSummary{}, nil)

	if err := testutils.RunApp(command, "l0 environment list"); err != nil {
		t.Fatal(err)
	}
}

func TestReadEnvironment(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Client.EXPECT().
		ReadEnvironment("env_id").
		Return(&models.Environment{}, nil)

	if err := testutils.RunApp(command, "l0 environment get env_name"); err != nil {
		t.Fatal(err)
	}
}

func TestReadEnvironmentInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	cases := map[string]string{
		"Missing NAME arg": "l0 environment get",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := testutils.RunApp(command, input); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestLinkEnvironments(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

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

	input := "l0 environment link --bi-directional=false env_name1 env_name2"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestLinkEnvironmentsBidirectional(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

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

	input := "l0 environment link env_name1 env_name2"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestLinkEnvironmentsInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	cases := map[string]string{
		"Missing SOURCE arg": "l0 environment link",
		"Missing DEST arg":   "l0 environment link src",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := testutils.RunApp(command, input); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestUnlinkEnvironments(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

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

	input := "l0 environment unlink --bi-directional=false env_name1 env_name2"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestUnlinkEnvironmentsBidirectional(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

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

	input := "l0 environment unlink env_name1 env_name2"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestUnlinkEnvironmentsInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	cases := map[string]string{
		"Missing SOURCE arg": "l0 environment unlink",
		"Missing DEST arg":   "l0 environment unlink src",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := testutils.RunApp(command, input); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}
