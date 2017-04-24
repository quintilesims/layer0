package command

import (
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Plan() cli.Command {
	return cli.Command{
		Name:      "plan",
		Usage:     "plan Layer0 instances",
		ArgsUsage: "NAME",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			instance := instance.NewInstance(args["NAME"])
			if err := instance.Plan(); err != nil {
				return err
			}

			return nil
		},
	}
}
