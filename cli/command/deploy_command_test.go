package command

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/urfave/cli"
	"testing"
)

func TestCreateDeploy(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(tc.Command())

	file, close := tempFile(t, "dockerrun")
	defer close()

	tc.Client.EXPECT().
		CreateDeploy("name", []byte("dockerrun")).
		Return(&models.Deploy{}, nil)

	c := getCLIContext(t, Args{file.Name(), "name"}, nil)
	if err := command.Create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDeploy_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing PATH arg": getCLIContext(t, nil, nil),
		"Missing NAME arg": getCLIContext(t, Args{"path"}, nil),
	}

	for name, c := range contexts {
		if err := command.Create(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestDeleteDeploy(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("deploy", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		DeleteDeploy("id").
		Return(nil)

	c := getCLIContext(t, Args{"name"}, nil)
	if err := command.Delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDeploy_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Delete(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestGetDeploy(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(tc.Command())

	tc.Resolver.EXPECT().
		Resolve("deploy", "name").
		Return([]string{"id"}, nil)

	tc.Client.EXPECT().
		GetDeploy("id").
		Return(&models.Deploy{}, nil)

	c := getCLIContext(t, Args{"name"}, nil)
	if err := command.Get(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetDeploy_userInputErrors(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(tc.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.Get(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestListDeploys(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(tc.Command())

	tc.Client.EXPECT().
		ListDeploys().
		Return([]*models.Deploy{}, nil)

	c := getCLIContext(t, nil, Flags{"all": true})
	if err := command.List(c); err != nil {
		t.Fatal(err)
	}
}

func TestFilterDeploys(t *testing.T) {
	input := []*models.Deploy{
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

	output, err := filterDeploys(input)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(output), 3)
	// max 'a', 'b', and 'c' deploys
	testutils.AssertInSlice(t, input[2], output)
	testutils.AssertInSlice(t, input[5], output)
	testutils.AssertInSlice(t, input[8], output)
}
