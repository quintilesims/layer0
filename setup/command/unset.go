package command

import (
	"fmt"

	"github.com/urfave/cli"
)

func (f *CommandFactory) Unset() cli.Command {
	return cli.Command{
		Name:      "unset",
		Usage:     "unset an input variable for a Layer0 instance's Terraform module",
		ArgsUsage: "NAME KEY",
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME", "KEY")
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			if err := instance.Unset(args["KEY"]); err != nil {
				return err
			}

			fmt.Println("Unset complete!")
			return nil
		},
	}
}
