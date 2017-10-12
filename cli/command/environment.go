package command

import (
	"fmt"

	"github.com/urfave/cli"
)

type EnvironmentCommand struct {
	*CommandBase
}

func NewEnvironmentCommand(b *CommandBase) *EnvironmentCommand {
	return &EnvironmentCommand{b}
}

func (e *EnvironmentCommand) Command() cli.Command {
	return cli.Command{
		Name: "environment",
		Action: func(c *cli.Context) error {
			id, err := e.resolveSingleEntityIDHelper(c.Args().Get(0), c.Args().Get(1))
			if err != nil {
				return err
			}

			fmt.Printf("Resolved: '%s'\n", id)
			return nil
		},
	}
}
