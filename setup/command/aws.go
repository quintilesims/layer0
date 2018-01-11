package command

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/setup/aws"
	"github.com/urfave/cli"
)

var awsFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "aws-access-key",
		Usage:  "access key portion of an AWS key",
		EnvVar: config.AWS_ACCESS_KEY_ID,
	},
	cli.StringFlag{
		Name:   "aws-secret-key",
		Usage:  "secret key portion on an AWS key",
		EnvVar: config.AWS_SECRET_ACCESS_KEY,
	},
}

func (f *CommandFactory) newAWSProviderHelper(c *cli.Context, region string) (*aws.Provider, error) {
	// first grab default config settings
	config := defaults.Get().Config

	// use static credentials if passed in by the user
	accessKey := c.String("aws-access-key")
	secretKey := c.String("aws-secret-key")
	if accessKey != "" && secretKey != "" {
		staticCreds := credentials.NewStaticCredentials(accessKey, secretKey, "")
		config.WithCredentials(staticCreds)
	} else {
		logrus.Debugf("aws-access-key or aws-secret-key was not specified. Using default credentials")
	}

	// ensure credentials are available
	if _, err := config.Credentials.Get(); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "NoCredentialProviders" {
			text := "No valid AWS credentials found. Please specify an AWS access key and secret key using "
			text += "their corresponding flags or environment variables."
			text += "`l0-setup init --aws-access-key <value> --aws-secret-key <value>`"
			return nil, fmt.Errorf(text)
		}

		return nil, err
	}

	// ensure that the correct region is set for AWS services
	// that have region-specific operations
	config.WithRegion(region)

	return f.NewAWSProvider(config), nil
}
