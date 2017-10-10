package command

import "github.com/urfave/cli"

type DeployCommand struct {
	*CommandMediator
}

func NewDeploy(m *CommandMediator) cli.Command {
	// a later reference d - commented out for now to avoid compile-time error
	// d = &DeployCommand{
	// 	CommandMediator: m,
	// }

	return cli.Command{
		Name:        "deploy",
		Usage:       "manage layer0 deploys",
		Subcommands: []cli.Command{},
	}
}
