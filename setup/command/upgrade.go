package command

import (
	"fmt"

	"github.com/urfave/cli"
)

func (f *CommandFactory) Upgrade() cli.Command {
	return cli.Command{
		Name:      "upgrade",
		Usage:     "upgrade a Layer0 instance to a new version",
		ArgsUsage: "NAME VERSION",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "force",
				Usage: "skips confirmation prompt",
			},
		},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME", "VERSION")
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			if err := instance.Upgrade(args["VERSION"], c.Bool("force")); err != nil {
				return err
			}

			fmt.Printf("Everything looks good! You are now ready to run 'l0-setup apply %s'\n", args["NAME"])
			return nil
		},
	}
}
