package command

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

func (f *CommandFactory) newAWSClientHelper(c *cli.Context) (*awsc.Client, error) {
	// use default credentials and region settings
	awsConfig := defaults.Get().Config

	// use static credentials if passed in by the user
	accessKey := c.String(config.FlagAWSAccessKey.GetName())
	secretKey := c.String(config.FlagAWSSecretKey.GetName())
	if accessKey != "" && secretKey != "" {
		staticCreds := credentials.NewStaticCredentials(accessKey, secretKey, "")
		awsConfig.WithCredentials(staticCreds)
	} else {
		log.Println("[DEBUG] aws-access-key or aws-secret-key was not specified. Using default credentials")
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

	region := c.String(config.FlagAWSRegion.GetName())
	awsConfig.WithRegion(region)

	sess := session.New(awsConfig)
	return f.NewAWSClient(sess), nil
}
