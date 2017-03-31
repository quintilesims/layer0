package command

import (
	"fmt"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Config() cli.Command {
	return cli.Command{
		Name:      "config",
		Usage:     "configure a layer0 instance",
		ArgsUsage: "NAME",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			instance, err := getInstance(f.InstanceFactory, c)
			if err != nil {
				return err
			}

			// todo: overwrites any variables going into the main.tf module via string flags
			print(instance)
			return fmt.Errorf("not implemented")
		},
	}
}
