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
	}
}

type APIConfig interface {
	Port() int
	AccessKey() string
	SecretKey() string
	Region() string
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
	vars := map[string]error{
		FLAG_AWS_ACCESS_KEY: fmt.Errorf("AWS Access Key not set! (EnvVar: %s)", ENVVAR_AWS_ACCESS_KEY),
		FLAG_AWS_SECRET_KEY: fmt.Errorf("AWS Secret Key not set! (EnvVar: %s)", ENVVAR_AWS_SECRET_KEY),
		FLAG_AWS_REGION:     fmt.Errorf("AWS Region not set! (EnvVar: %s)", ENVVAR_AWS_REGION),
	}

	for name, err := range vars {
		if c.C.String(name) == "" {
			return err
		}
	}

	return nil
}

func (c *ContextAPIConfig) Port() int {
        return c.C.Int(FLAG_PORT)
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
