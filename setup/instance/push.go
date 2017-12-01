package instance

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

func (l *LocalInstance) Push(s s3iface.S3API) error {
	if err := l.assertExists(); err != nil {
		return err
	}

	bucket, err := l.Output(OUTPUT_S3_BUCKET)
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

		key := strings.Replace(path, l.Dir, "terraform", 1)
		log.Printf("[INFO] Pushing %s to s3://%s/%s", path, bucket, key)

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

	return filepath.Walk(l.Dir, pushFiles)
}
