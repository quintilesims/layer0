package command

import (
	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Init() cli.Command {
	return cli.Command{
		Name:      "init",
		Usage:     "initialize a new layer0 instance",
		ArgsUsage: "NAME",
		Flags: []cli.Flag{
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
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			accessKey := c.String("access-key")
			if accessKey == "" {
				return fmt.Errorf("Please provide an aws access key using the associated flag or environment variable")
			}

			secretKey := c.String("secret-key")
			if secretKey == "" {
				return fmt.Errorf("Please provide an aws secret key using the associated flag or environment variable")
			}

			region := c.String("region")
			if region == "" {
				return fmt.Errorf("Please provide an aws region using the associated flag or environment variable")
			}

			keyPair := c.String("key-pair")
			if keyPair == "" {
				return fmt.Errorf("Please provide an aws key pair using the associated flag or environment variable")
			}

			inst, err := f.InstanceFactory.NewInstance(args["NAME"])
			if err != nil {
				return err
			}

			exists, err := inst.Exists()
			if err != nil {
				return err
			}

			if exists {
				return fmt.Errorf("Instance '%s' already exists", inst.Name)
			}

			config := instance.InstanceConfig{
				AccessKey:        accessKey,
				SecretKey:        secretKey,
				Region:           region,
				KeyPair:          keyPair,
				DockerConfigPath: c.String("docker-config"),
			}

			if err := inst.Init(config); err != nil {
				return err
			}

			fmt.Printf("Successfully created new instance '%s'\n", inst.Name())
			return nil
		},
	}
}
