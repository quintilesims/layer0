package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/quintilesims/layer0/common/aws/provider"
	"io/ioutil"
	"os"
)

type Provider interface {
	PutObject(string, string, []byte) error
	ListObjects(string, string) ([]string, error)
	GetObject(string, string) ([]byte, error)
	DeleteObject(string, string) error
	PutObjectFromFile(bucket, key, path string) error
	GetObjectToFile(bucket, key, path string, filemode os.FileMode) error
}

type s3Internal interface {
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	DeleteObject(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	ListObjects(*s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
}

type S3 struct {
	credProvider provider.CredProvider
	region       string
	Connect      func() (s3Internal, error)
}

func Connect(credProvider provider.CredProvider, region string) (s3Internal, error) {
	connection, err := provider.GetS3Connection(credProvider, region)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func NewS3(credProvider provider.CredProvider, region string) (Provider, error) {
	s3 := S3{
		credProvider: credProvider,
		region:       region,
		Connect:      func() (s3Internal, error) { return Connect(credProvider, region) },
	}

	_, err := s3.Connect()
	if err != nil {
		return nil, err
	}

	return &s3, nil
}

func (this *S3) PutObject(bucket, key string, body []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.PutObject(input)
	if err != nil {
		return err
	}

	return nil
}

func (this *S3) PutObjectFromFile(bucket, key, path string) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return this.PutObject(bucket, key, body)
}

func (this *S3) GetObject(bucket, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	resp, err := connection.GetObject(input)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.Bytes(), nil
}

func (this *S3) DeleteObject(bucket, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteObject(input)
	if err != nil {
		return err
	}

	return nil
}

func (this *S3) GetObjectToFile(bucket, key, path string, fileMode os.FileMode) error {
	body, err := this.GetObject(bucket, key)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, body, fileMode)
}

func (this *S3) ListObjects(bucket, prefix string) ([]string, error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	resp, err := connection.ListObjects(input)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, object := range resp.Contents {
		keys = append(keys, *object.Key)
	}

	return keys, nil
}
