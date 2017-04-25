package command

import (
	"github.com/urfave/cli"
)

func (f *CommandFactory) Plan() cli.Command {
	return cli.Command{
		Name:      "plan",
		Usage:     "show the planned operation(s) to run during the next 'apply'",
		ArgsUsage: "NAME",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			if err := instance.Plan(); err != nil {
				return err
			}

			return nil
		},
	}
}
