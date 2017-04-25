package instance

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/docker/docker/pkg/homedir"
	"io/ioutil"
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
	instanceBuckets, err := listInstanceBuckets(s)
	if err != nil {
		return nil, err
	}

	instances := []string{}
	for instance := range instanceBuckets {
		instances = append(instances, instance)
	}

	return instances, nil
}
