package config

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func flagsToEnvVar() map[string]string {
	return map[string]string{
		FLAG_PORT:                      ENVVAR_PORT,
		FLAG_TOKEN:                     ENVVAR_TOKEN,
		FLAG_LOCK_EXPIRY:               ENVVAR_LOCK_EXPIRY,
		FLAG_DEBUG:                     ENVVAR_DEBUG,
		FLAG_INSTANCE:                  ENVVAR_INSTANCE,
		FLAG_AWS_ACCOUNT_ID:            ENVVAR_AWS_ACCOUNT_ID,
		FLAG_AWS_ACCESS_KEY:            ENVVAR_AWS_ACCESS_KEY,
		FLAG_AWS_SECRET_KEY:            ENVVAR_AWS_SECRET_KEY,
		FLAG_AWS_REGION:                ENVVAR_AWS_REGION,
		FLAG_AWS_VPC:                   ENVVAR_AWS_VPC,
		FLAG_AWS_LINUX_AMI:             ENVVAR_AWS_LINUX_AMI,
		FLAG_AWS_S3_BUCKET:             ENVVAR_AWS_S3_BUCKET,
		FLAG_AWS_INSTANCE_PROFILE:      ENVVAR_AWS_INSTANCE_PROFILE,
		FLAG_AWS_DYNAMO_TAG_TABLE:      ENVVAR_AWS_DYNAMO_TAG_TABLE,
		FLAG_AWS_DYNAMO_LOCK_TABLE:     ENVVAR_AWS_DYNAMO_LOCK_TABLE,
		FLAG_AWS_PRIVATE_SUBNETS:       ENVVAR_AWS_PRIVATE_SUBNETS,
		FLAG_AWS_PUBLIC_SUBNETS:        ENVVAR_AWS_PUBLIC_SUBNETS,
		FLAG_AWS_LOG_GROUP_NAME:        ENVVAR_AWS_LOG_GROUP_NAME,
		FLAG_AWS_TIME_BETWEEN_REQUESTS: ENVVAR_AWS_TIME_BETWEEN_REQUESTS,
		FLAG_AWS_SSH_KEY_PAIR:          ENVVAR_AWS_SSH_KEY_PAIR,
		FLAG_AWS_MAX_RETRIES:           ENVVAR_AWS_MAX_RETRIES,
	}
}

func APIFlags() []cli.Flag {
	flags := flagsToEnvVar()

	return []cli.Flag{
		cli.IntFlag{
			Name:   FLAG_PORT,
			Value:  DefaultPort,
			EnvVar: flags[FLAG_PORT],
		},
		cli.StringFlag{
			Name:   FLAG_TOKEN,
			EnvVar: flags[FLAG_TOKEN],
		},
		cli.DurationFlag{
			Name:   FLAG_LOCK_EXPIRY,
			Value:  DefaultLockExpiry,
			EnvVar: flags[FLAG_LOCK_EXPIRY],
		},
		cli.BoolFlag{
			Name:   FLAG_DEBUG,
			EnvVar: flags[FLAG_DEBUG],
		},
		cli.StringFlag{
			Name:   FLAG_INSTANCE,
			EnvVar: flags[FLAG_INSTANCE],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_ACCOUNT_ID,
			EnvVar: flags[FLAG_AWS_ACCOUNT_ID],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_ACCESS_KEY,
			EnvVar: flags[FLAG_AWS_ACCESS_KEY],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_SECRET_KEY,
			EnvVar: flags[FLAG_AWS_SECRET_KEY],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_REGION,
			Value:  DefaultAWSRegion,
			EnvVar: flags[FLAG_AWS_REGION],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_VPC,
			EnvVar: flags[FLAG_AWS_VPC],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_LINUX_AMI,
			EnvVar: flags[FLAG_AWS_LINUX_AMI],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_S3_BUCKET,
			EnvVar: flags[FLAG_AWS_S3_BUCKET],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_INSTANCE_PROFILE,
			EnvVar: flags[FLAG_AWS_INSTANCE_PROFILE],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_DYNAMO_TAG_TABLE,
			EnvVar: flags[FLAG_AWS_DYNAMO_TAG_TABLE],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_DYNAMO_LOCK_TABLE,
			EnvVar: flags[FLAG_AWS_DYNAMO_LOCK_TABLE],
		},
		cli.StringSliceFlag{
			Name:   FLAG_AWS_PUBLIC_SUBNETS,
			EnvVar: flags[FLAG_AWS_PUBLIC_SUBNETS],
		},
		cli.StringSliceFlag{
			Name:   FLAG_AWS_PRIVATE_SUBNETS,
			EnvVar: flags[FLAG_AWS_PRIVATE_SUBNETS],
		},
		cli.StringFlag{
			Name:   FLAG_AWS_LOG_GROUP_NAME,
			EnvVar: flags[FLAG_AWS_LOG_GROUP_NAME],
		},
		cli.DurationFlag{
			Name:   FLAG_AWS_TIME_BETWEEN_REQUESTS,
			Value:  DefaultTimeBetweenRequests,
			EnvVar: flags[FLAG_AWS_TIME_BETWEEN_REQUESTS],
			Usage:  "duration [h,m,s,ms,ns]",
		},
		cli.StringFlag{
			Name:   FLAG_AWS_SSH_KEY_PAIR,
			EnvVar: flags[FLAG_AWS_SSH_KEY_PAIR],
		},
		cli.IntFlag{
			Name:   FLAG_AWS_MAX_RETRIES,
			Value:  50,
			EnvVar: flags[FLAG_AWS_MAX_RETRIES],
		},
	}
}

type APIConfig interface {
	Port() int
	ParseAuthToken() (string, string, error)
	AccountID() string
	AccessKey() string
	SecretKey() string
	Region() string
	Instance() string
	VPC() string
	LinuxAMI() string
	S3Bucket() string
	InstanceProfile() string
	PublicSubnets() []string
	PrivateSubnets() []string
	DynamoTagTable() string
	DynamoLockTable() string
	LogGroupName() string
	SSHKeyPair() string
	LockExpiry() time.Duration
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
	flags := flagsToEnvVar()

	requiredVars := []string{
		FLAG_INSTANCE,
		FLAG_TOKEN,
		FLAG_AWS_ACCOUNT_ID,
		FLAG_AWS_ACCESS_KEY,
		FLAG_AWS_SECRET_KEY,
		FLAG_AWS_VPC,
		FLAG_AWS_LINUX_AMI,
		FLAG_AWS_S3_BUCKET,
		FLAG_AWS_INSTANCE_PROFILE,
		FLAG_AWS_DYNAMO_TAG_TABLE,
		FLAG_AWS_DYNAMO_LOCK_TABLE,
		FLAG_AWS_PUBLIC_SUBNETS,
		FLAG_AWS_PRIVATE_SUBNETS,
		FLAG_AWS_LOG_GROUP_NAME,
		FLAG_AWS_SSH_KEY_PAIR,
	}

	for _, name := range requiredVars {
		if !c.C.IsSet(name) {
			return fmt.Errorf("required flag '%s' or environment variable '%s' is not set", name, flags[name])
		}
	}

	return nil
}

func (c *ContextAPIConfig) Port() int {
	return c.C.Int(FLAG_PORT)
}

func (c *ContextAPIConfig) AuthToken() string {
	return c.C.String(FLAG_TOKEN)
}

func (c *ContextAPIConfig) ParseAuthToken() (string, string, error) {
	token, err := base64.StdEncoding.DecodeString(c.AuthToken())
	if err != nil {
		return "", "", fmt.Errorf("Auth Token is not in valid base64 format: %v", err)
	}

	split := strings.Split(string(token), ":")
	if len(split) != 2 {
		return "", "", fmt.Errorf("Auth Token must be in format 'user:pass' and base64 encoded")
	}

	return split[0], split[1], nil
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

func (c *ContextAPIConfig) SSHKeyPair() string {
	return c.C.String(FLAG_AWS_SSH_KEY_PAIR)
}

func (c *ContextAPIConfig) LinuxAMI() string {
	return c.C.String(FLAG_AWS_LINUX_AMI)
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

func (c *ContextAPIConfig) DynamoLockTable() string {
	return c.C.String(FLAG_AWS_DYNAMO_LOCK_TABLE)
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

func (c *ContextAPIConfig) LockExpiry() time.Duration {
	return c.C.Duration(FLAG_LOCK_EXPIRY)
}

func (c *ContextAPIConfig) TimeBetweenRequests() time.Duration {
	return c.C.Duration(FLAG_AWS_TIME_BETWEEN_REQUESTS)
}

func (c *ContextAPIConfig) MaxRetries() int {
	return c.C.Int(FLAG_AWS_MAX_RETRIES)
}
