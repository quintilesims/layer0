package command

import (
	"fmt"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/setup/aws"
	"github.com/urfave/cli"
)

var awsFlags = []cli.Flag{
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
}

func (f *CommandFactory) newAWSProviderHelper(c *cli.Context) (*aws.Provider, error) {
	accessKey := c.String("aws-access-key")
	if accessKey == "" {
		return nil, fmt.Errorf("AWS Access Key not set! (EnvVar: %s)", config.AWS_ACCESS_KEY_ID)
	}

	secretKey := c.String("aws-secret-key")
	if secretKey == "" {
		return nil, fmt.Errorf("AWS Secret Key not set! (EnvVar: %s)", config.AWS_SECRET_ACCESS_KEY)
	}

	region := c.String("aws-region")
	if region == "" {
		return nil, fmt.Errorf("AWS Region not set! (EnvVar: %s)", config.AWS_REGION)
	}

	return f.NewAWSProvider(accessKey, secretKey, region), nil
}
