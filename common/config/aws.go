package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func GetAWSConfig(creds *credentials.Credentials, region string) *aws.Config {
	awsConfig := &aws.Config{}
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(region)
	awsConfig.WithMaxRetries(DEFAULT_MAX_RETRIES)

	return awsConfig
}
