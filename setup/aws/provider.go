package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

const DEFAULT_AWS_REGION = "us-west-2"

type Provider struct {
	EC2 ec2iface.EC2API
	S3  s3iface.S3API
}

func NewProvider(config *aws.Config) *Provider {
	session := session.New(config)
	return &Provider{
		EC2: ec2.New(session),
		S3:  s3.New(session),
	}
}
