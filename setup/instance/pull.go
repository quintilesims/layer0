package instance

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"io/ioutil"
	"os"
	"strings"
)

func (l *LocalInstance) Pull(s s3iface.S3API) error {
	buckets, err := listLocalInstanceBuckets(s)
	if err != nil {
		return err
	}

	bucket, ok := buckets[l.Name]
	if !ok {
		return fmt.Errorf("S3 bucket for instance '%s' does not exist!", l.Name)
	}

	// get all of the files in the 'terraform' directory of the bucket
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("terraform"),
	}

	output, err := s.ListObjects(input)
	if err != nil {
		return err
	}

	// create the parent directories if they don't already exist
	if err := os.MkdirAll(l.Dir, 0700); err != nil {
		return err
	}

	for _, content := range output.Contents {
		path := strings.Replace(aws.StringValue(content.Key), "terraform", l.Dir, 1)
		logrus.Infof("Pulling s3://%s/%s to %s\n", bucket, aws.StringValue(content.Key), path)

		input := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    content.Key,
		}

		output, err := s.GetObject(input)
		if err != nil {
			return err
		}

		// if we are looking at a directory, create it locally (don't do any file io)
		if aws.StringValue(output.ContentType) == "application/x-directory" {
			if err := os.MkdirAll(path, 0700); err != nil {
				return err
			}

			continue
		}

		// otherwise, write the file locally
		data, err := ioutil.ReadAll(output.Body)
		if err != nil {
			return err
		}
		defer output.Body.Close()

		if err := ioutil.WriteFile(path, data, 0644); err != nil {
			return err
		}
	}

	return nil
}
