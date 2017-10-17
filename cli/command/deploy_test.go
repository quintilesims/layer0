package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCreateDeploy(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.Command())

	file, close := createTempFile(t, "dpl_file")
	defer close()

	req := models.CreateDeployRequest{
		DeployName: "dpl_name",
		DeployFile: []byte("dpl_file"),
	}

	base.Client.EXPECT().
		CreateDeploy(req).
		Return("jid", nil)

	job := &models.Job{
		Status: "Completed",
		Result: "dpl_id",
	}

	base.Client.EXPECT().
		ReadJob("jid").
		Return(job, nil)

	base.Client.EXPECT().
		ReadDeploy("dpl_id").
		Return(&models.Deploy{}, nil)

	c := getCLIContext(t, []string{file.Name(), "dpl_name"}, nil)
	if err := command.create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDeploy_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing PATH arg": getCLIContext(t, nil, nil),
		"Missing NAME arg": getCLIContext(t, []string{"path"}, nil),
	}

	for name, c := range contexts {
		if err := command.create(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestDeleteDeploy(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"id"}, nil)

	base.Client.EXPECT().
		DeleteDeploy("id").
		Return("jid", nil)

	job := &models.Job{
		Status: "Completed",
		Result: "dpl_id",
	}

	base.Client.EXPECT().
		ReadJob("jid").
		Return(job, nil)

	c := getCLIContext(t, []string{"dpl_name"}, nil)
	if err := command.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDeploy_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.delete(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestReadDeploy(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("deploy", "dpl_name").
		Return([]string{"id"}, nil)

	base.Client.EXPECT().
		ReadDeploy("id").
		Return(&models.Deploy{}, nil)

	c := getCLIContext(t, []string{"dpl_name"}, nil)
	if err := command.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadDeploy_userInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": getCLIContext(t, nil, nil),
	}

	for name, c := range contexts {
		if err := command.read(c); err == nil {
			t.Fatalf("%s: error was nil!", name)
		}
	}
}

func TestListDeploys(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewDeployCommand(base.Command())

	base.Client.EXPECT().
		ListDeploys().
		Return([]*models.DeploySummary{}, nil)

	c := getCLIContext(t, nil, map[string]interface{}{"all": true})
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestFilterDeploySummaries(t *testing.T) {
	input := []*models.DeploySummary{
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
