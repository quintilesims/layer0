package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type Client struct {
	AutoScaling autoscalingiface.AutoScalingAPI
	EC2         ec2iface.EC2API
	ECS         ecsiface.ECSAPI
	ELB         elbiface.ELBAPI
	IAM         iamiface.IAMAPI
	S3          s3iface.S3API
}

func NewClient(config *aws.Config) *Client {
	session := session.New(config)
	return &Client{
		AutoScaling: autoscaling.New(session),
		EC2:         ec2.New(session),
		ECS:         ecs.New(session),
		ELB:         elb.New(session),
		IAM:         iam.New(session),
		S3:          s3.New(session),
	}
}
