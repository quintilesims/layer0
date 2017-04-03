package command

import (
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Apply() cli.Command {
	return cli.Command{
		Name:      "apply",
		Usage:     "apply a layer0 instance",
		ArgsUsage: "NAME",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			instance := instance.NewInstance(args["NAME"])
			return instance.Apply()
		},
	}
}
