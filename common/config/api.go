package config

import (
	"fmt"

	"github.com/urfave/cli"
)

func APIFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			// todo: renamed from 'LAYER0_API_PORT'
			Name:   FLAG_PORT,
			Value:  9090,
			EnvVar: ENVVAR_PORT,
		},
		cli.BoolFlag{
			// todo: renamed from 'LAYER0_LOG_LEVEL'
			Name:   FLAG_DEBUG,
			EnvVar: ENVVAR_DEBUG,
		},
		cli.StringFlag{
			Name:   FLAG_INSTANCE,
			EnvVar: ENVVAR_INSTANCE,
		},
		cli.StringFlag{
			Name:   FLAG_AWS_ACCOUNT_ID,
			EnvVar: ENVVAR_AWS_ACCOUNT_ID,
		},
		cli.StringFlag{
			Name:   FLAG_AWS_ACCESS_KEY,
			EnvVar: ENVVAR_AWS_ACCESS_KEY,
		},
		cli.StringFlag{
			Name:   FLAG_AWS_SECRET_KEY,
			EnvVar: ENVVAR_AWS_SECRET_KEY,
		},
		cli.StringFlag{
			Name:   FLAG_AWS_REGION,
			Value:  "us-west-2",
			EnvVar: ENVVAR_AWS_REGION,
		},
		cli.StringFlag{
			Name:   FLAG_AWS_VPC,
			EnvVar: ENVVAR_AWS_VPC,
		},
		cli.StringFlag{
			Name:   FLAG_AWS_LINUX_AMI,
			EnvVar: ENVVAR_AWS_LINUX_AMI,
		},
		cli.StringFlag{
			Name:   FLAG_AWS_WINDOWS_AMI,
			EnvVar: ENVVAR_AWS_WINDOWS_AMI,
		},
		cli.StringFlag{
			Name:   FLAG_AWS_S3_BUCKET,
			EnvVar: ENVVAR_AWS_S3_BUCKET,
		},
		cli.StringFlag{
			Name:   FLAG_AWS_INSTANCE_PROFILE,
			EnvVar: ENVVAR_AWS_INSTANCE_PROFILE,
		},
		cli.StringFlag{
			Name:   FLAG_AWS_DYNAMO_TAG_TABLE,
			EnvVar: ENVVAR_AWS_DYNAMO_TAG_TABLE,
		},
		cli.StringSliceFlag{
			Name:   FLAG_AWS_PUBLIC_SUBNETS,
			EnvVar: ENVVAR_AWS_PUBLIC_SUBNETS,
		},
		cli.StringSliceFlag{
			Name:   FLAG_AWS_PRIVATE_SUBNETS,
			EnvVar: ENVVAR_AWS_PRIVATE_SUBNETS,
		},
	}
}

type APIConfig interface {
	Port() int
	AccountID() string
	AccessKey() string
	SecretKey() string
	Region() string
	Instance() string
	VPC() string
	LinuxAMI() string
	WindowsAMI() string
	S3Bucket() string
	InstanceProfile() string
	PublicSubnets() []string
	PrivateSubnets() []string
	DynamoTagTable() string
}

type ContextAPIConfig struct {
	C *cli.Context
}

func NewContextAPIConfig(c *cli.Context) *ContextAPIConfig {
	return &ContextAPIConfig{
		C: c,
	}
}

func (c *ContextAPIConfig) Validate() error {
	stringVars := map[string]error{
		FLAG_INSTANCE:             fmt.Errorf("Layer0 Instance not set! (EnvVar: %s)", ENVVAR_INSTANCE),
		FLAG_AWS_ACCOUNT_ID:       fmt.Errorf("AWS Account ID not set! (EnvVar: %s)", ENVVAR_AWS_ACCOUNT_ID),
		FLAG_AWS_ACCESS_KEY:       fmt.Errorf("AWS Access Key not set! (EnvVar: %s)", ENVVAR_AWS_ACCESS_KEY),
		FLAG_AWS_SECRET_KEY:       fmt.Errorf("AWS Secret Key not set! (EnvVar: %s)", ENVVAR_AWS_SECRET_KEY),
		FLAG_AWS_REGION:           fmt.Errorf("AWS Region not set! (EnvVar: %s)", ENVVAR_AWS_REGION),
		FLAG_AWS_VPC:              fmt.Errorf("AWS VPC not set! (EnvVar: %s)", ENVVAR_AWS_VPC),
		FLAG_AWS_LINUX_AMI:        fmt.Errorf("AWS Linux AMI not set! (EnvVar: %s)", ENVVAR_AWS_LINUX_AMI),
		FLAG_AWS_WINDOWS_AMI:      fmt.Errorf("AWS Windows AMI not set! (EnvVar: %s)", ENVVAR_AWS_WINDOWS_AMI),
		FLAG_AWS_S3_BUCKET:        fmt.Errorf("AWS S3 Bucket not set! (EnvVar: %s)", ENVVAR_AWS_S3_BUCKET),
		FLAG_AWS_INSTANCE_PROFILE: fmt.Errorf("AWS Instance Profile not set! (EnvVar: %s)", ENVVAR_AWS_INSTANCE_PROFILE),
		FLAG_AWS_DYNAMO_TAG_TABLE: fmt.Errorf("AWS Dynamo Tag Table not set! (EnvVar: %s)", ENVVAR_AWS_DYNAMO_TAG_TABLE),
	}

	for name, err := range stringVars {
		if c.C.String(name) == "" {
			return err
		}
	}

	stringSliceVars := map[string]error{
		FLAG_AWS_PUBLIC_SUBNETS:  fmt.Errorf("Layer0 Public Subnets not set! (EnvVar: %s)", ENVVAR_AWS_PUBLIC_SUBNETS),
		FLAG_AWS_PRIVATE_SUBNETS: fmt.Errorf("Layer0 Private Subnets not set! (EnvVar: %s)", ENVVAR_AWS_PRIVATE_SUBNETS),
	}

	for name, err := range stringSliceVars {
		if len(c.C.StringSlice(name)) == 0 {
			return err
		}
	}

	return nil
}

func (c *ContextAPIConfig) Port() int {
	return c.C.Int(FLAG_PORT)
}

func (c *ContextAPIConfig) Instance() string {
	return c.C.String(FLAG_INSTANCE)
}

func (c *ContextAPIConfig) AccountID() string {
	return c.C.String(FLAG_AWS_ACCOUNT_ID)
}

func (c *ContextAPIConfig) AccessKey() string {
	return c.C.String(FLAG_AWS_ACCESS_KEY)
}

func (c *ContextAPIConfig) SecretKey() string {
	return c.C.String(FLAG_AWS_SECRET_KEY)
}

func (c *ContextAPIConfig) Region() string {
	return c.C.String(FLAG_AWS_REGION)
}

func (c *ContextAPIConfig) VPC() string {
	return c.C.String(FLAG_AWS_VPC)
}

func (c *ContextAPIConfig) LinuxAMI() string {
	return c.C.String(FLAG_AWS_LINUX_AMI)
}

func (c *ContextAPIConfig) WindowsAMI() string {
	return c.C.String(FLAG_AWS_WINDOWS_AMI)
}

func (c *ContextAPIConfig) S3Bucket() string {
	return c.C.String(FLAG_AWS_S3_BUCKET)
}

func (c *ContextAPIConfig) InstanceProfile() string {
	return c.C.String(FLAG_AWS_INSTANCE_PROFILE)
}

func (c *ContextAPIConfig) DynamoTagTable() string {
	return c.C.String(FLAG_AWS_DYNAMO_TAG_TABLE)
}

func (c *ContextAPIConfig) PublicSubnets() []string {
	return c.C.StringSlice(FLAG_AWS_PUBLIC_SUBNETS)
}

func (c *ContextAPIConfig) PrivateSubnets() []string {
	return c.C.StringSlice(FLAG_AWS_PRIVATE_SUBNETS)
}
