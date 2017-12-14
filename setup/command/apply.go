package command

import (
	"fmt"

	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Apply() cli.Command {
	return cli.Command{
		Name:      "apply",
		Usage:     "create and/or update resources for a Layer0 instance",
		ArgsUsage: "NAME",
		Flags: []cli.Flag{
			config.FlagAWSAccessKey,
			config.FlagAWSSecretKey,
			config.FlagAWSRegion,
			cli.BoolFlag{
				Name:  "quick",
				Usage: "skips verification checks that normally run after 'terraform apply' has completed",
			},
			cli.BoolTFlag{
				Name:  "push",
				Usage: "setting it to false skips pushing local tfstate to s3 (default: true)",
			},
		},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			if err := instance.Apply(!c.Bool("quick")); err != nil {
				return err
			}

			if c.Bool("push") {
				provider, err := f.newAWSClientHelper(c)
				if err != nil {
					return err
				}

				if err := instance.Push(provider.S3); err != nil {
					return err
				}
			}

			fmt.Println("Apply complete!")
			return nil
		},
	}
}
