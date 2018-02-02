package command

import (
	"fmt"
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeploy(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.CommandBase()).Command()

	file, delete := testutils.TempFile(t, "content")
	defer delete()

	req := models.CreateDeployRequest{
		DeployName: "dpl_name",
		DeployFile: []byte("content"),
	}

	base.Client.EXPECT().
		CreateDeploy(req).
		Return("dpl_id", nil)

	base.Client.EXPECT().
		ReadDeploy("dpl_id").
		Return(&models.Deploy{}, nil)

	input := fmt.Sprintf("l0 deploy create %s dpl_name", file.Name())
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDeployInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.CommandBase()).Command()

	cases := map[string]string{
		"Missing PATH arg":    "l0 deploy create",
		"Missing NAME arg":    "l0 deploy create path",
		"PATH does not exist": "l0 deploy create badpath name",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := testutils.RunApp(command, input); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestDeleteDeploy(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	base.Client.EXPECT().
		DeleteDeploy("dpl_id").
		Return(nil)

	if err := testutils.RunApp(command, "l0 deploy delete dpl_name"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDeployInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.CommandBase()).Command()

	cases := map[string]string{
		"Missing NAME arg": "l0 deploy delete",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := testutils.RunApp(command, input); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestReadDeploy(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.CommandBase()).Command()

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"dpl_id"}, nil)

	base.Client.EXPECT().
		ReadDeploy("dpl_id").
		Return(&models.Deploy{}, nil)

	if err := testutils.RunApp(command, "l0 deploy get dpl_name"); err != nil {
		t.Fatal(err)
	}
}

func TestReadDeployInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.CommandBase()).Command()

	cases := map[string]string{
		"Missing NAME arg": "l0 deploy get",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := testutils.RunApp(command, input); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestListDeploys(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.CommandBase()).Command()

	base.Client.EXPECT().
		ListDeploys().
		Return([]models.DeploySummary{}, nil)

	if err := testutils.RunApp(command, "l0 deploy list"); err != nil {
		t.Fatal(err)
	}
}

func TestFilterDeploySummaries(t *testing.T) {
	input := []models.DeploySummary{
		{DeployName: "a", DeployID: "a.1", Version: "1"},
		{DeployName: "a", DeployID: "a.2", Version: "2"},
		{DeployName: "a", DeployID: "a.3", Version: "3"},
		{DeployName: "b", DeployID: "b.1", Version: "1"},
		{DeployName: "b", DeployID: "b.2", Version: "2"},
		{DeployName: "b", DeployID: "b.3", Version: "3"},
		{DeployName: "c", DeployID: "c.9", Version: "9"},
		{DeployName: "c", DeployID: "c.10", Version: "10"},
		{DeployName: "c", DeployID: "c.11", Version: "11"},
		{DeployID: "nameless.1", Version: "1"},
		{DeployID: "nameless.2", Version: "2"},
	}

	output, err := filterDeploySummaries(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(output), 3)
	// max 'a', 'b', and 'c' deploys
	assert.Contains(t, output, input[2])
	assert.Contains(t, output, input[5])
	assert.Contains(t, output, input[8])
}
