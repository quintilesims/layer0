package testutils

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/quintilesims/layer0/common/config"
)

const (
	ENVVAR_TEST_AWS_DYNAMO_TAG_TABLE  = "LAYER0_TEST_AWS_DYNAMO_TAG_TABLE"
	ENVVAR_TEST_AWS_DYNAMO_LOCK_TABLE = "LAYER0_TEST_AWS_DYNAMO_LOCK_TABLE"
)

func GetTestAWSSession() *session.Session {
	accessKey := os.Getenv(config.ENVVAR_AWS_ACCESS_KEY)
	secretKey := os.Getenv(config.ENVVAR_AWS_SECRET_KEY)
	region := os.Getenv(config.ENVVAR_AWS_REGION)
	if region == "" {
		region = config.DefaultAWSRegion
	}

	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}

	return session.New(awsConfig)
}
