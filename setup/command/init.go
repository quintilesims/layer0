package command

import (
	"fmt"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

// todo: allow docker config path flag
func (f *CommandFactory) Init() cli.Command {
	// !! do not use defaults for the these flags, otherwise the inputs
	// will *always* be overwritten by the default value

	return cli.Command{
		Name:  "init",
		Usage: "initialize or reconfigure a layer0 instance",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "module-source",
				Usage:  "path to Layer0 module",
				EnvVar: "LAYER0_MODULE_SOURCE",
			},
			cli.StringFlag{
				Name:   "version",
				Usage:  "version of Layer0 to use",
				EnvVar: "LAYER0_VERSION",
			},
			cli.StringFlag{
				Name:   "aws-access-key",
				Usage:  "AWS access key id",
				EnvVar: config.AWS_ACCESS_KEY_ID,
			},
			cli.StringFlag{
				Name:   "aws-secret-key",
				Usage:  "AWS secret access key",
				EnvVar: config.AWS_SECRET_ACCESS_KEY,
			},
			cli.StringFlag{
				Name:   "aws-region",
				Usage:  "AWS region",
				EnvVar: config.AWS_REGION,
			},
			cli.StringFlag{
				Name:   "aws-key-pair",
				Usage:  "AWS key pair",
				EnvVar: config.AWS_KEY_PAIR,
			},
		},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			overrides := map[string]interface{}{}
			if v := c.String("module-source"); v != "" {
				overrides[instance.INPUT_SOURCE] = v
			}

			if v := c.String("aws-access-key"); v != "" {
				overrides[instance.INPUT_AWS_ACCESS_KEY] = v
			}

			if v := c.String("aws-secret-key"); v != "" {
				overrides[instance.INPUT_AWS_SECRET_KEY] = v
			}

			if v := c.String("aws-region"); v != "" {
				overrides[instance.INPUT_AWS_REGION] = v
			}

			if v := c.String("aws-key-pair"); v != "" {
				overrides[instance.INPUT_AWS_KEY_PAIR] = v
			}

			instance := instance.NewInstance(args["NAME"])
			if err := instance.Init(c, overrides); err != nil {
				return err
			}

			fmt.Println("Initialization complete!")
			return nil
		},
	}
}
