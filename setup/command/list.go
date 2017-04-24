package command

import (
	"fmt"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) List() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list local layer0 instances",
		Flags: []cli.Flag{
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
		},
		Action: func(c *cli.Context) error {
			accessKey := c.String("aws-access-key")
			if accessKey == "" {
				return fmt.Errorf("AWS Access Key not set! (EnvVar: %s)", config.AWS_ACCESS_KEY_ID)
			}

			secretKey := c.String("aws-secret-key")
			if secretKey == "" {
				return fmt.Errorf("AWS Secret Key not set! (EnvVar: %s)", config.AWS_SECRET_ACCESS_KEY)
			}

			s3 := newS3(accessKey, secretKey)
			remote, err := instance.ListRemoteInstances(s3)
			if err != nil {
				return err
			}

			local, err := instance.ListLocalInstances()
			if err != nil {
				return err
			}

			catalog := map[string]string{}
			for _, instance := range local {
				catalog[instance] += "l"
			}

			for _, instance := range remote {
				catalog[instance] += "r"
			}

			for instance, token := range catalog {
				fmt.Printf("%2s    %s\n", token, instance)
			}

			return nil
		},
	}
}
