package command

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func (f *CommandFactory) Set() cli.Command {
	return cli.Command{
		Name:      "set",
		Usage:     "set input variable(s) for a Layer0 instance's Terraform module",
		ArgsUsage: "NAME",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "input",
				Usage: "Specify an input using key=val format",
			},
		},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			inputs := map[string]interface{}{}
			for _, input := range c.StringSlice("input") {
				split := strings.Split(input, "=")
				if len(split) != 2 {
					return fmt.Errorf("Invalid input format '%s'", input)
				}

				inputs[split[0]] = split[1]
			}

			instance := f.NewInstance(args["NAME"])
			if err := instance.Set(inputs); err != nil {
				return err
			}

			fmt.Println("Set complete!")
			return nil
		},
	}
}
