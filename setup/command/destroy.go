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
			instance, err := getInstance(f.InstanceFactory, c)
			if err != nil {
				return err
			}

			if err := instance.Destroy(); err != nil {
				return err
			}

			fmt.Printf("Successfully destroyed '%s'\n", instance.Name())
			return nil
		},
	}
}
