package instance

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/docker/docker/pkg/homedir"
	"io/ioutil"
	"strings"
)

func ListLocalInstances() ([]string, error) {
	dir := fmt.Sprintf("%s/.layer0", homedir.Get())
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	instances := []string{}
	for _, file := range files {
		if file.IsDir() {
			instances = append(instances, file.Name())
		}
	}

	return instances, nil
}

func ListRemoteInstances(s *s3.S3) ([]string, error) {
	buckets, err := listLayer0Buckets(s)
	if err != nil {
		return nil, err
	}

	instances := []string{}
	for _, bucket := range buckets {
		// layer0 bucket name format: 'layer0-<instance>-<account_id>'
		if split := strings.Split(bucket, "-"); len(split) == 3 {
			instances = append(instances, split[1])
		}
	}

	return instances, nil
}

func listLayer0Buckets(s *s3.S3) ([]string, error) {
	output, err := s.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	buckets := []string{}
	for _, bucket := range output.Buckets {
		name := aws.StringValue(bucket.Name)
		if strings.HasPrefix(name, "layer0-") {
			buckets = append(buckets, name)
		}
	}

	return buckets, nil
}
