package command

import (
	"fmt"

	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Pull() cli.Command {
	return cli.Command{
		Name:      "pull",
		Usage:     "pull a Layer0 instance configuration from S3",
		ArgsUsage: "NAME",
		Flags:     awsFlags,
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			i := f.NewInstance(args["NAME"])
			region, err := i.Output(instance.OUTPUT_AWS_REGION)
			if err != nil {
				return err
			}

			provider, err := f.newAWSProviderHelper(c, region)
			if err != nil {
				return err
			}

			if err := i.Pull(provider.S3); err != nil {
				return err
			}

			fmt.Println("Pull complete!")
			return nil
		},
	}
}
