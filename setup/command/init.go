package command

import (
	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/layer0/setup/layer0"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Init() cli.Command {
	return cli.Command{
		Name:  "init",
		Usage: "initialize a new layer0 instance",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "name",
				Usage:  "name for your new layer0 instance",
				EnvVar: "LAYER0_INSTANCE_NAME",
			},
			cli.StringFlag{
				Name:   "access-key",
				Usage:  "aws access key id",
				EnvVar: "AWS_ACCESS_KEY_ID",
			},
			cli.StringFlag{
				Name:   "secret-key",
				Usage:  "aws secret access key",
				EnvVar: "AWS_SECRET_ACCESS_KEY",
			},
			cli.StringFlag{
				Name:   "region",
				Usage:  "aws region",
				Value:  "us-west-2",
				EnvVar: "AWS_REGION",
			},
			cli.StringFlag{
				Name:   "key-pair",
				Usage:  "aws key pair",
				EnvVar: "AWS_KEY_PAIR",
			},
			cli.StringFlag{
				Name:   "docker-config",
				Usage:  "path to your docker config file",
				Value:  fmt.Sprintf("%s/.docker/config.json", homedir.Get()),
				EnvVar: "DOCKER_CONFIG_PATH",
			},
		},
		Action: func(c *cli.Context) error {
			config := layer0.InstanceConfig{
				Name:      c.String("name"),
				AccessKey: c.String("access_key"),
				SecretKey: c.String("secret_key"),
			}

			if _, err := f.context.Init(config); err != nil {
				return err
			}

			return nil
		},
	}
}
