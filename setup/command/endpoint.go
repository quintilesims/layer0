package command

import (
	"fmt"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Endpoint() cli.Command {
	return cli.Command{
		Name:      "endpoint",
		Usage:     "endpoint layer0 instances",
		ArgsUsage: "NAME",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			instance, err := getInstance(f.InstanceFactory, c)
			if err != nil {
				return err
			}

			print(instance)
			return fmt.Errorf("not implemented")
		},
	}
}
