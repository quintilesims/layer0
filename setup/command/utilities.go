package command

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

var s3Flags = []cli.Flag{
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

func extractArgs(received []string, names ...string) (map[string]string, error) {
	args := map[string]string{}
	for i, name := range names {
		if len(received)-1 < i {
			return nil, fmt.Errorf("Argument %s is required", name)
		}

		args[name] = received[i]
	}

	return args, nil
}

func newS3(c *cli.Context) (*s3.S3, error) {
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

	session := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(region),
	})

	return s3.New(session), nil
}
