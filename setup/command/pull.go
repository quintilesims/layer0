package command

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	aws_provider "github.com/quintilesims/layer0/setup/aws"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Pull() cli.Command {
	return cli.Command{
		Name:      "pull",
		Usage:     "pull a Layer0 instance configuration from S3",
		ArgsUsage: "NAME",
		Flags:     awsFlags,
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			// Use the default AWS region first to retrieve the list of buckets
			provider, err := f.newAWSProviderHelper(c, aws_provider.DEFAULT_AWS_REGION)
			if err != nil {
				return err
			}

			remoteInstanceBucket, err := getRemoteInstanceBucket(provider.S3, args["NAME"])
			if err != nil {
				return err
			}

			region, err := getBucketLocation(provider.S3, remoteInstanceBucket)
			if err != nil {
				return err
			}

			if region == "" {
				// See: https://docs.aws.amazon.com/AmazonS3/latest/API/RESTBucketGETlocation.html
				// When GetBucketLocation returns an empty string, the bucket is in us-east-1
				region = "us-east-1"
			}

			// Change the AWS provider configuration to match the region of the bucket
			// to pull from
			provider, err = f.newAWSProviderHelper(c, region)
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			if err := instance.Pull(provider.S3); err != nil {
				return err
			}

			fmt.Println("Pull complete!")
			return nil
		},
	}
}

func getRemoteInstanceBucket(s s3iface.S3API, instanceName string) (string, error) {
	listBucketsOutput, err := s.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return "", err
	}

	for _, bucket := range listBucketsOutput.Buckets {
		bucketName := aws.StringValue(bucket.Name)

		if split := strings.Split(bucketName, "-"); len(split) == 3 && split[0] == "layer0" && split[1] == instanceName {
			return bucketName, nil
		}
	}

	return "", fmt.Errorf("No S3 bucket found for given instance name")
}

func getBucketLocation(s s3iface.S3API, bucketName string) (string, error) {
	getBucketLocationInput := &s3.GetBucketLocationInput{}
	getBucketLocationInput.SetBucket(bucketName)

	getBucketLocationOutput, err := s.GetBucketLocation(getBucketLocationInput)
	if err != nil {
		return "", err
	}

	return aws.StringValue(getBucketLocationOutput.LocationConstraint), nil
}
