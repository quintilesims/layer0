package command

import "github.com/urfave/cli"

type DeployCommand struct {
	*CommandMediator
}

func NewDeployCommand(m *CommandMediator) *DeployCommand {
	return &DeployCommand{
		CommandMediator: m,
	}
}

func (d *DeployCommand) Command() cli.Command {
	return cli.Command{
		Name:        "deploy",
		Usage:       "manage layer0 deploys",
		Subcommands: []cli.Command{},
	}
}
