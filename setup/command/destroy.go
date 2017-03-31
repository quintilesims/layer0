package command

import (
	"fmt"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Destroy() cli.Command {
	return cli.Command{
		Name:      "destroy",
		Usage:     "destroy a layer0 instance",
		ArgsUsage: "NAME",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			return fmt.Errorf("not implemented")
		},
	}
}
