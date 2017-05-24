package instance

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/docker/docker/pkg/homedir"
	"io/ioutil"
	"strings"
	"os"
)

func ListLocalInstances() ([]string, error) {
	dir := fmt.Sprintf("%s/.layer0", homedir.Get())
	files, err := ioutil.ReadDir(dir)
	if err != nil && !os.IsNotExist(err) {
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

func ListRemoteInstances(s s3iface.S3API) ([]string, error) {
	instanceBuckets, err := listRemoteInstanceBuckets(s)
	if err != nil {
		return nil, err
	}

	instances := []string{}
	for instance := range instanceBuckets {
		instances = append(instances, instance)
	}

	return instances, nil
}

func listRemoteInstanceBuckets(s s3iface.S3API) (map[string]string, error) {
	output, err := s.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	instanceBuckets := map[string]string{}
	for _, bucket := range output.Buckets {
		name := aws.StringValue(bucket.Name)

		// layer0 bucket name format: 'layer0-<instance>-<account_id>'
		if split := strings.Split(name, "-"); len(split) == 3 && split[0] == "layer0" {
			instanceBuckets[split[1]] = name
		}
	}

	return instanceBuckets, nil
}
