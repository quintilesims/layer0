package command

import (
	"fmt"
	"net/url"
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
		EnvironmentType:  "static",
		InstanceType:     "t2.small",
		Scale:            2,
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
	input += "--scale 2 "
	input += "--os linux "
	input += "--ami ami123 "
	input += fmt.Sprintf("--user-data %s ", file.Name())
	input += "env_name "

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestCreateEnvironmentInputErrors(t *testing.T) {
	testInputErrors(t, NewEnvironmentCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 environment create",
	})
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

	input := "l0 environment delete env_name"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentDeleteRecursive(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	loadBalancerToDelete := []models.LoadBalancerSummary{
		{
			LoadBalancerID: "lb_id",
			EnvironmentID:  "env_id",
		},
	}

	base.Client.EXPECT().
		ListLoadBalancers().
		Return(loadBalancerToDelete, nil)

	base.Client.EXPECT().
		DeleteLoadBalancer("lb_id").
		Return(nil)

	taskToDelete := []models.TaskSummary{
		{
			TaskID:        "tsk_id",
			EnvironmentID: "env_id",
		},
	}

	base.Client.EXPECT().
		ListTasks().
		Return(taskToDelete, nil)

	base.Client.EXPECT().
		DeleteTask("tsk_id").
		Return(nil)

	serviceToDelete := []models.ServiceSummary{
		{
			ServiceID:     "svc_id",
			EnvironmentID: "env_id",
		},
	}

	base.Client.EXPECT().
		ListServices().
		Return(serviceToDelete, nil)

	base.Client.EXPECT().
		DeleteService("svc_id").
		Return(nil)

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	base.Client.EXPECT().
		DeleteEnvironment("env_id").
		Return(nil)

	input := "l0 environment delete "
	input += "--recursive "
	input += "env_name"

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEnvironmentInputErrors(t *testing.T) {
	testInputErrors(t, NewEnvironmentCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 environment delete",
	})
}

func TestListEnvironments(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	base.Client.EXPECT().
		ListEnvironments().
		Return([]models.EnvironmentSummary{}, nil)

	input := "l0 environment list"
	if err := testutils.RunApp(command, input); err != nil {
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

	input := "l0 environment get env_name"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestReadEnvironmentInputErrors(t *testing.T) {
	testInputErrors(t, NewEnvironmentCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 environment get",
	})
}

func TestReadEnvironmentLogs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewEnvironmentCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("environment", "env_name").
		Return([]string{"env_id"}, nil)

	query := url.Values{
		"tail":  []string{"100"},
		"start": []string{"start"},
		"end":   []string{"end"},
	}

	base.Client.EXPECT().
		ReadEnvironmentLogs("env_id", query).
		Return([]models.LogFile{}, nil)

	input := "l0 environment logs "
	input += "--tail 100 "
	input += "--start start "
	input += "--end end "
	input += "env_name"

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestReadEnvironmentLogsInputErrors(t *testing.T) {
	testInputErrors(t, NewEnvironmentCommand(nil).Command(), map[string]string{
		"Missing NAME arg": "l0 environment logs",
	})
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

	input := "l0 environment link "
	input += "--bi-directional=false "
	input += "env_name1 env_name2"

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
	testInputErrors(t, NewEnvironmentCommand(nil).Command(), map[string]string{
		"Missing SOURCE arg": "l0 environment link",
		"Missing DEST arg":   "l0 environment link src",
	})
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

	input := "l0 environment unlink "
	input += "--bi-directional=false "
	input += "env_name1 env_name2"
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
	testInputErrors(t, NewEnvironmentCommand(nil).Command(), map[string]string{
		"Missing SOURCE arg": "l0 environment unlink",
		"Missing DEST arg":   "l0 environment unlink src",
	})
}
