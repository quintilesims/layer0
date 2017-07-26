package command

import (
	"fmt"
	"strings"

	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Init() cli.Command {
	return cli.Command{
		Name:  "init",
		Usage: "initialize or reconfigure a Layer0 instance",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "docker-path",
				Usage: "path to docker config.json file",
			},
			cli.StringFlag{
				Name:  "docker-creds-helper-path",
				Usage: "path to dokcer credential helper",
			},
			cli.StringFlag{
				Name:  "module-source",
				Usage: instance.INPUT_SOURCE_DESCRIPTION,
			},
			cli.StringFlag{
				Name:  "version",
				Usage: instance.INPUT_VERSION_DESCRIPTION,
			},
			cli.StringFlag{
				Name:  "aws-access-key",
				Usage: instance.INPUT_AWS_ACCESS_KEY_DESCRIPTION,
			},
			cli.StringFlag{
				Name:  "aws-secret-key",
				Usage: instance.INPUT_AWS_SECRET_KEY_DESCRIPTION,
			},
			cli.StringFlag{
				Name:  "aws-region",
				Usage: instance.INPUT_AWS_REGION_DESCRIPTION,
			},
			cli.StringFlag{
				Name:  "aws-ssh-key-pair",
				Usage: instance.INPUT_AWS_SSH_KEY_PAIR_DESCRIPTION,
			},
		},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			// do not use defaults for the override flags, otherwise the inputs
			// will *always* be overwritten by the default value
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

			if v := c.String("aws-ssh-key-pair"); v != "" {
				overrides[instance.INPUT_AWS_SSH_KEY_PAIR] = v
			}

			dockerPath := strings.Replace(c.String("docker-path"), "~", homedir.Get(), -1)
			dockerCredsHelperPath := strings.Replace(c.String("docker-creds-helper-path"), "~", homedir.Get(), -1)

			instance := f.NewInstance(args["NAME"])
			if err := instance.Init(dockerPath, dockerCredsHelperPath, overrides); err != nil {
				return err
			}

			fmt.Println("Initialization complete!")
			return nil
		},
	}
}
