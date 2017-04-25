package command

import (
	"fmt"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Push() cli.Command {
	return cli.Command{
		Name:      "push",
		Usage:     "todo",
		ArgsUsage: "NAME",
		Flags:     s3Flags,
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			s3, err := newS3(c)
			if err != nil {
				return err
			}

			instance := instance.NewInstance(args["NAME"])
			if err := instance.Push(s3); err != nil {
				return err
			}

			fmt.Println("Push complete!")
			return nil
		},
	}
}
