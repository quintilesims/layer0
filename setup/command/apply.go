package command

import (
	"fmt"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Apply() cli.Command {
	return cli.Command{
		Name:      "apply",
		Usage:     "apply a layer0 instance",
		ArgsUsage: "NAME",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			instance, err := getInstance(f.InstanceFactory, c)
			if err != nil {
				return err
			}

			if err := instance.Apply(); err != nil {
				return err
			}

			fmt.Printf("Successfully applied instance '%s'\n", instance.Name())
			return nil
		},
	}
}
