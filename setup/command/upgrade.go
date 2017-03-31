package command

import (
	"fmt"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Upgrade() cli.Command {
	return cli.Command{
		Name:      "upgrade",
		Usage:     "upgrade Layer0 instances",
		ArgsUsage: "NAME",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			return fmt.Errorf("not implemented")
		},
	}
}
