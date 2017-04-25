package instance

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
)

func listInstanceBuckets(s *s3.S3) (map[string]string, error) {
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
