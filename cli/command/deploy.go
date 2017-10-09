package command

import "github.com/urfave/cli"

func (f *CommandFactory) Deploy() cli.Command {
	return cli.Command{
		Name:        "deploy",
		Usage:       "manage layer0 deploys",
		Subcommands: []cli.Command{},
	}
}
