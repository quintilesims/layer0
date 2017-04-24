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
	output, err := s.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	instances := []string{}
	for _, bucket := range output.Buckets {
		name := aws.StringValue(bucket.Name)

		// layer0 bucket name format: 'layer0-<instance>-<account id>'
		split := strings.Split(name, "-")
		if len(split) == 3 && split[0] == "layer0" {
			instances = append(instances, split[1])
		}
	}

	return instances, nil
}
