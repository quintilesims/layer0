package command

import (
	"github.com/urfave/cli"
)

func (f *CommandFactory) List() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list local layer0 instances",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
