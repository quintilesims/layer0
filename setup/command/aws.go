package command

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

var awsFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "aws-access-key",
		Usage:  "access key portion of an AWS key",
		EnvVar: config.ENVVAR_AWS_ACCESS_KEY,
	},
	cli.StringFlag{
		Name:   "aws-secret-key",
		Usage:  "secret key portion on an AWS key",
		EnvVar: config.ENVVAR_AWS_SECRET_KEY,
	},
	cli.StringFlag{
		Name:   "aws-region",
		Usage:  "AWS region",
		EnvVar: config.ENVVAR_AWS_REGION,
	},
}

func (f *CommandFactory) newAWSClientHelper(c *cli.Context) (*awsc.Client, error) {
	// use default credentials and region settings
	awsConfig := defaults.Get().Config

	// use static credentials if passed in by the user
	accessKey := c.String("aws-access-key")
	secretKey := c.String("aws-secret-key")
	if accessKey != "" && secretKey != "" {
		staticCreds := credentials.NewStaticCredentials(accessKey, secretKey, "")
		awsConfig.WithCredentials(staticCreds)
	} else {
		logrus.Debugf("aws-access-key or aws-secret-key was not specified. Using default credentials")
	}

	// ensure credentials are available
	if _, err := awsConfig.Credentials.Get(); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "NoCredentialProviders" {
			text := "No valid AWS credentials found. Please specify an AWS access key and secret key using "
			text += "their corresponding flags or environment variables"
			return nil, fmt.Errorf(text)
		}

		return nil, err
	}

	// use region if passed in by the user
	awsConfig.WithRegion(config.DefaultAWSRegion)
	if region := c.String("aws-region"); region != "" {
		awsConfig.WithRegion(region)
	} else {
		logrus.Debugf("aws-region was not specified. Using default")
	}

	return f.NewAWSClient(awsConfig), nil
}
