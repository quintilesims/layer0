package command

import (
	"fmt"

	setup_instance "github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Apply() cli.Command {
	return cli.Command{
		Name:      "apply",
		Usage:     "create and/or update resources for a Layer0 instance",
		ArgsUsage: "NAME",
		Flags: append(awsFlags, []cli.Flag{
			cli.BoolFlag{
				Name:  "quick",
				Usage: "skips verification checks that normally run after 'terraform apply' has completed",
			},
			cli.BoolTFlag{
				Name:  "push",
				Usage: "setting it to false skips pushing local tfstate to s3 (default: true)",
			},
		}...),
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
			}

			fmt.Println("Apply complete!")
			return nil
		},
	}
}
