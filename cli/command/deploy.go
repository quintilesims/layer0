package command

import "github.com/urfave/cli"

type DeployCommand struct {
	*CommandBase
}

func NewDeployCommand(b *CommandBase) *DeployCommand {
	return &DeployCommand{b}
}

func (d *DeployCommand) Command() cli.Command {
	return cli.Command{
		Name:        "deploy",
		Usage:       "manage layer0 deploys",
		Subcommands: []cli.Command{},
	}
}
