package command

import (
	"fmt"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Apply() cli.Command {
	return cli.Command{
		Name:      "apply",
		Usage:     "create and/or update resources for your layer0 instance",
		ArgsUsage: "NAME",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			if err := instance.Apply(); err != nil {
				return err
			}

			fmt.Println("Apply complete!")
			return nil
		},
	}
}
