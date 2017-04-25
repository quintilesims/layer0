package instance

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/layer0/setup/terraform"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Instance struct {
	Name      string
	Dir       string
	Terraform *terraform.Terraform
}

func NewInstance(name string) *Instance {
	dir := fmt.Sprintf("%s/.layer0/%s", homedir.Get(), name)

	return &Instance{
		Name: name,
		Dir:  dir,
	}
}

func (i *Instance) Apply() error {
	if err := i.Terraform.Apply(i.Dir); err != nil {
		return err
	}

	endpoint, err := i.Output(OUTPUT_ENDPOINT)
	if err != nil {
		return err
	}

	return i.waitForHealthyAPI(endpoint, time.Minute*10)
}

func (i *Instance) waitForHealthyAPI(endpoint string, timeout time.Duration) error {
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second * 15) {
		log.Printf("Waiting for API Service to be healthy... (%v)\n", time.Since(start))
		resp, err := http.Get(endpoint)
		if err != nil {
			log.Println("Error getting api: ", err)
			continue
		}

		defer resp.Body.Close()
		if code := resp.StatusCode; code < 200 || code > 299 {
			log.Println("API returned non-200 status code: %d", code)
			continue
		}

		return nil
	}

	return fmt.Errorf("API Service was not healthy after %v", timeout)
}

func (i *Instance) Destroy(force bool) error {
	// todo: use layer0 client to destroy all resources

	if err := i.Terraform.Destroy(i.Dir, force); err != nil {
		return err
	}

	return os.RemoveAll(i.Dir)
}

func (i *Instance) Init(c *cli.Context, inputOverrides map[string]interface{}) error {
	if err := os.MkdirAll(i.Dir, 0700); err != nil {
		return err
	}

	// load terraform config from ~/.layer0/<instance>/main.tf.json, or create a new one
	config, err := i.loadMainConfig()
	if err != nil {
		return err
	}

	// add/update the inputs of the terraform config
	if err := i.setMainModuleInputs(config, inputOverrides); err != nil {
		return err
	}

	// save the terraform config as ~/.layer0/<instance>/main.tf.json
	path := fmt.Sprintf("%s/main.tf.json", i.Dir)
	if err := terraform.WriteConfig(path, config); err != nil {
		return err
	}

	// create/write ~/.layer0/<instance>/outputs.tf.json
	output := &terraform.Config{
		Outputs: MainModuleOutputs,
	}

	outPath := fmt.Sprintf("%s/outputs.tf.json", i.Dir)
	if err := terraform.WriteConfig(outPath, output); err != nil {
		return err
	}

	// run `terraform get` to download terraform modules
	if err := i.Terraform.Get(i.Dir); err != nil {
		return err
	}

	// run `terraform fmt` to validate the terraform syntax
	if err := i.Terraform.FMT(i.Dir); err != nil {
		return err
	}

	return nil
}

func (i *Instance) loadMainConfig() (*terraform.Config, error) {
	path := fmt.Sprintf("%s/main.tf.json", i.Dir)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return terraform.LoadConfig(path)
	}

	return terraform.NewConfig(), nil
}

func (i *Instance) setMainModuleInputs(config *terraform.Config, inputOverrides map[string]interface{}) error {
	// create the 'main' module if it doesn't already exist
	if _, ok := config.Modules["main"]; !ok {
		config.Modules["main"] = terraform.Module{}
	}

	module := config.Modules["main"]
	for _, input := range MainModuleInputs {
		// if the user specified a cli flag or env var, use that for the input
		if v, ok := inputOverrides[input.Name]; ok {
			module[input.Name] = v
			continue
		}

		// prompt the user for a new/updated input
		v, err := input.Prompt(module[input.Name])
		if err != nil {
			return err
		}

		module[input.Name] = v
	}

	// the 'name' input is always the name of the layer0 instance
	module["name"] = i.Name
	return nil
}

func (i *Instance) Output(key string) (string, error) {
	return i.Terraform.Output(i.Dir, key)
}

func (i *Instance) Plan() error {
	return i.Terraform.Plan(i.Dir)
}

func (i *Instance) Push(s *s3.S3) error {
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

func (i *Instance) Pull(s *s3.S3) error {
	if err := os.MkdirAll(i.Dir, 0700); err != nil {
		return err
	}

	bucket, err := i.getBucket(s)
	if err != nil {
		return err
	}

	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("terraform"),
	}

	output, err := s.ListObjects(input)
	if err != nil {
		return err
	}

	for _, content := range output.Contents {
		path := strings.Replace(aws.StringValue(input.Key), "terraform", i.Dir, 1)
		log.Printf("Pulling s3://%s/%s to %s\n", bucket, aws.StringValue(content.Key), path)

		input := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    content.Key,
		}

		output, err := s.GetObject(input)
		if err != nil {
			return err
		}

		if aws.StringValue(output.ContentType) == "application/x-directory" {
			if err := os.MkdirAll(path, 0700); err != nil {
				return err
			}

			continue
		}

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

func (i *Instance) getBucket(s *s3.S3) (string, error) {
	buckets, err := listLayer0Buckets(s)
	if err != nil {
		return "", err
	}

	for _, bucket := range buckets {
		// layer0 bucket name format: 'layer0-<instance>-<account_id>'
		split := strings.Split(bucket, "-")
		if len(split) == 3 && split[1] == i.Name {
			return bucket, nil
		}
	}

	return "", fmt.Errorf("S3 bucket for instance '%s' does not exist!", i.Name)
}
