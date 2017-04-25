package instance

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (i *Instance) Push(s *s3.S3) error {
	if err := i.assertExists(); err != nil {
		return err
	}

	bucket, err := i.Output(OUTPUT_S3_BUCKET)
	if err != nil {
		return err
	}

	pushFiles := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() || strings.Contains(path, ".terraform") {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		key := strings.Replace(path, i.Dir, "terraform", 1)
		log.Printf("Pushing %s to s3://%s/%s\n", path, bucket, key)

		input := &s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			Body:        bytes.NewReader(data),
			ContentType: aws.String("text/json"),
		}

		if _, err := s.PutObject(input); err != nil {
			return err
		}

		return nil
	}

	return filepath.Walk(i.Dir, pushFiles)
}
