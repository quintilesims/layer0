package provider

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/quintilesims/layer0/common/config"
)

const (
	AP_NORTHEAST_1 = "ap-northeast-1"
	AP_SOUTHEAST_1 = "ap-southeast-1"
	AP_SOUTHEAST_2 = "ap-southeast-2"
	EU_CENTRAL_1   = "eu-central-1"
	EU_WEST_1      = "eu-west-1"
	SA_EAST_1      = "sa-east-1"
	US_EAST_1      = "us-east-1"
	US_WEST_1      = "us-west-1"
	US_WEST_2      = "us-west-2"
)

var regionIsValid = func(region string) (isValid bool) {
	regions := []string{
		AP_NORTHEAST_1,
		AP_SOUTHEAST_1,
		AP_SOUTHEAST_2,
		EU_CENTRAL_1,
		EU_WEST_1,
		SA_EAST_1,
		US_EAST_1,
		US_WEST_1,
		US_WEST_2,
	}

	for _, r := range regions {
		if region == r {
			isValid = true
			break
		}
	}

	return
}

var Ticker *time.Ticker

func init() {
	delay, err := time.ParseDuration(config.AWSTimeBetweenRequests())
	if err != nil {
		return
	}

	Ticker = time.NewTicker(delay)
}

var getConfig = func(credProvider CredProvider, region string) (sess *session.Session, err error) {
	fmt.Printf("getConfig called with %v and %v\n", credProvider, region)
	fmt.Printf("ticker: %v\n", Ticker)
	if !regionIsValid(region) {
		err = fmt.Errorf("Region '%s' is not a valid region!", region)
		return
	}

	access_key, err := credProvider.GetAWSAccessKeyID()
	if err != nil {
		return
	}

	secret_key, err := credProvider.GetAWSSecretAccessKey()
	if err != nil {
		return
	}

	creds := credentials.NewStaticCredentials(access_key, secret_key, "")
	sess = session.New(config.GetAWSConfig(creds, config.AWSRegion()))
	sess.Handlers.Send.PushBack(func(r *request.Request) {
		<-Ticker.C
	})

	return
}

var GetElasticBeanstalkConnection = func(credProvider CredProvider, region string) (connection *elasticbeanstalk.ElasticBeanstalk, err error) {
	sess, err := getConfig(credProvider, region)
	if err != nil {
		return
	}

	connection = elasticbeanstalk.New(sess)
	return
}

var GetCloudFormationConnection = func(credProvider CredProvider, region string) (connection *cloudformation.CloudFormation, err error) {
	sess, err := getConfig(credProvider, region)
	if err != nil {
		return
	}

	connection = cloudformation.New(sess)
	return
}

var GetCloudWatchLogsConnection = func(credProvider CredProvider, region string) (connection *cloudwatchlogs.CloudWatchLogs, err error) {
	sess, err := getConfig(credProvider, region)
	if err != nil {
		return
	}

	connection = cloudwatchlogs.New(sess)
	return
}

var GetEC2Connection = func(credProvider CredProvider, region string) (connection *ec2.EC2, err error) {
	sess, err := getConfig(credProvider, region)
	if err != nil {
		return
	}

	connection = ec2.New(sess)
	return
}

var GetCloudWatchConnection = func(credProvider CredProvider, region string) (connection *cloudwatch.CloudWatch, err error) {
	sess, err := getConfig(credProvider, region)
	if err != nil {
		return
	}

	connection = cloudwatch.New(sess)
	return
}

var GetS3Connection = func(credProvider CredProvider, region string) (connection *s3.S3, err error) {
	sess, err := getConfig(credProvider, region)
	if err != nil {
		return
	}

	connection = s3.New(sess)
	return
}

var GetIAMConnection = func(credProvider CredProvider, region string) (connection *iam.IAM, err error) {
	sess, err := getConfig(credProvider, region)
	if err != nil {
		return
	}

	connection = iam.New(sess)
	return
}

var GetECSConnection = func(credProvider CredProvider, region string) (connection *ecs.ECS, err error) {
	sess, err := getConfig(credProvider, region)
	if err != nil {
		return
	}

	connection = ecs.New(sess)
	return
}

var GetELBConnection = func(credProvider CredProvider, region string) (connection *elb.ELB, err error) {
	sess, err := getConfig(credProvider, region)
	if err != nil {
		return
	}

	connection = elb.New(sess)
	return
}

var GetAutoScalingConnection = func(credProvider CredProvider, region string) (connection *autoscaling.AutoScaling, err error) {
	sess, err := getConfig(credProvider, region)
	if err != nil {
		return
	}

	connection = autoscaling.New(sess)
	return
}
