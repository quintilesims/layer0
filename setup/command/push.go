package command

import (
	"fmt"

	setup_instance "github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Push() cli.Command {
	return cli.Command{
		Name:      "push",
		Usage:     "push a Layer0 instance configuration to S3",
		ArgsUsage: "NAME",
		Flags:     awsFlags,
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			region, err := instance.Output(setup_instance.OUTPUT_AWS_REGION)
			if err != nil {
				return err
			}

			provider, err := f.newAWSProviderHelper(c, region)
			if err != nil {
				return err
			}

			if err := instance.Push(provider.S3); err != nil {
				return err
			}

			fmt.Println("Push complete!")
			return nil
		},
	}
}
