package command

import (
	"fmt"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Destroy() cli.Command {
	return cli.Command{
		Name:      "destroy",
		Usage:     "destroy all resources associated with your layer0 instance",
		ArgsUsage: "NAME",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "force",
			},
		},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			if err := instance.Destroy(c.Bool("force")); err != nil {
				return err
			}

			fmt.Println("Destroy complete!")
			return nil
		},
	}
}
