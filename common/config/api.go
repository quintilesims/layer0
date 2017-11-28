package config

import (
	"fmt"
	"time"

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
			Value:  DefaultAWSRegion,
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
			Name:   FLAG_AWS_DYNAMO_JOB_TABLE,
			EnvVar: ENVVAR_AWS_DYNAMO_JOB_TABLE,
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
		cli.StringFlag{
			Name:   FLAG_AWS_LOG_GROUP_NAME,
			EnvVar: ENVVAR_AWS_LOG_GROUP_NAME,
		},
		cli.DurationFlag{
			Name:   FLAG_AWS_TIME_BETWEEN_REQUESTS,
			Value:  10 * time.Millisecond,
			EnvVar: ENVVAR_AWS_TIME_BETWEEN_REQUESTS,
			Usage:  "duration [ns,us (or µs),ms,s,m,h]",
		},
		cli.IntFlag{
			Name:   FLAG_AWS_MAX_RETRIES,
			Value:  50,
			EnvVar: ENVVAR_AWS_MAX_RETRIES,
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
	DynamoJobTable() string
	DynamoTagTable() string
	LogGroupName() string
	TimeBetweenRequests() time.Duration
	MaxRetries() int
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
	requiredVars := []string{
		FLAG_INSTANCE,
		FLAG_AWS_ACCOUNT_ID,
		FLAG_AWS_ACCESS_KEY,
		FLAG_AWS_SECRET_KEY,
		FLAG_AWS_VPC,
		FLAG_AWS_LINUX_AMI,
		FLAG_AWS_WINDOWS_AMI,
		FLAG_AWS_S3_BUCKET,
		FLAG_AWS_INSTANCE_PROFILE,
		FLAG_AWS_DYNAMO_JOB_TABLE,
		FLAG_AWS_DYNAMO_TAG_TABLE,
		FLAG_AWS_PUBLIC_SUBNETS,
		FLAG_AWS_PRIVATE_SUBNETS,
		FLAG_AWS_LOG_GROUP_NAME,
	}

	for _, name := range requiredVars {
		if !c.C.IsSet(name) {
			return fmt.Errorf("Required Variable %s is not set!", name)
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

func (c *ContextAPIConfig) DynamoJobTable() string {
	return c.C.String(FLAG_AWS_DYNAMO_JOB_TABLE)
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

func (c *ContextAPIConfig) LogGroupName() string {
	return c.C.String(FLAG_AWS_LOG_GROUP_NAME)
}

func (c *ContextAPIConfig) TimeBetweenRequests() time.Duration {
	return c.C.Duration(FLAG_AWS_TIME_BETWEEN_REQUESTS)
}

func (c *ContextAPIConfig) MaxRetries() int {
	return c.C.Int(FLAG_AWS_MAX_RETRIES)
}
