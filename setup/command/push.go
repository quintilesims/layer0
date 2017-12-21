package command

import (
	"fmt"

	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Push() cli.Command {
	return cli.Command{
		Name:      "push",
		Usage:     "push a Layer0 instance configuration to S3",
		ArgsUsage: "NAME",
		Flags: []cli.Flag{
			config.FlagAWSAccessKey,
			config.FlagAWSSecretKey,
			config.FlagAWSRegion,
		},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			provider, err := f.newAWSClientHelper(c)
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			if err := instance.Push(provider.S3); err != nil {
				return err
			}

			fmt.Println("Push complete!")
			return nil
		},
	}
}
