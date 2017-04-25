package command

import (
	"fmt"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Pull() cli.Command {
	return cli.Command{
		Name:      "pull",
		Usage:     "todo",
		ArgsUsage: "NAME",
		Flags:     awsFlags,
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			provider, err := f.newAWSProviderHelper(c)
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			if err := instance.Pull(provider.S3); err != nil {
				return err
			}

			fmt.Println("Pull complete!")
			return nil
		},
	}
}
