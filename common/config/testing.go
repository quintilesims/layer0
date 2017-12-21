package config

import (
	"flag"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/urfave/cli"
)

type Option func(t *testing.T, c *cli.Context, flagSet *flag.FlagSet)

func NewTestContext(t *testing.T, args []string, flags map[string]interface{}, options ...Option) *cli.Context {
	flagSet := &flag.FlagSet{}
	c := cli.NewContext(&cli.App{}, flagSet, nil)

	for key, val := range flags {
		switch v := val.(type) {
		case bool:
			flagSet.Bool(key, v, "")
			c.Set(key, strconv.FormatBool(v))
		case int:
			flagSet.Int(key, v, "")
			c.Set(key, strconv.Itoa(v))
		case string:
			flagSet.String(key, v, "")
			c.Set(key, v)
		case []string:
			slice := cli.StringSlice(v)
			flagSet.Var(&slice, key, "")
		default:
			t.Fatalf("Unexpected flag type for '%s'", key)
		}
	}

	// add default global flags
	flagSet.String(FlagOutput.GetName(), DefaultOutput, "")
	flagSet.String(FlagTimeout.GetName(), DefaultTimeout.String(), "")
	flagSet.Bool(FlagNoWait.GetName(), false, "")

	for _, option := range options {
		option(t, c, flagSet)
	}

	if err := flagSet.Parse(args); err != nil {
		t.Fatal(err)
	}

	return c
}

func SetGlobalFlag(key, val string) Option {
	return func(t *testing.T, c *cli.Context, flagSet *flag.FlagSet) {
		if err := c.GlobalSet(key, val); err != nil {
			t.Fatal(err)
		}
	}
}

func SetNoWait(b bool) Option {
	return SetGlobalFlag(FlagNoWait.GetName(), strconv.FormatBool(b))
}

func SetTimeout(d time.Duration) Option {
	return SetGlobalFlag(FlagTimeout.GetName(), d.String())
}

func SetVersion(v string) Option {
	return func(t *testing.T, c *cli.Context, flagSet *flag.FlagSet) {
		c.App.Version = v
	}
}

func GetTestAWSSession() *session.Session {
	accessKey := os.Getenv(FlagAWSAccessKey.EnvVar)
	secretKey := os.Getenv(FlagAWSSecretKey.EnvVar)
	region := os.Getenv(FlagAWSRegion.EnvVar)
	if region == "" {
		region = DefaultAWSRegion
	}

	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}

	return session.New(awsConfig)
}
