package config

import (
	"fmt"
	"time"

	"github.com/urfave/cli"
)

func CLIFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   FLAG_OUTPUT,
			EnvVar: FLAG_OUTPUT,
			Value:  "text",
			Usage:  "output format [text,json]",
		},
		cli.DurationFlag{
			Name:   FLAG_TIMEOUT,
			EnvVar: ENVVAR_TIMEOUT,
			Value:  time.Minute * 15,
			Usage:  "timeout [h,m,s,ms]",
		},
		cli.BoolFlag{
			Name:   FLAG_DEBUG,
			EnvVar: ENVVAR_DEBUG,
			Usage:  "show debug output",
		},
		cli.StringFlag{
			Name:   FLAG_ENDPOINT,
			EnvVar: ENVVAR_ENDPOINT,
			Value:  "http://localhost:9090/",
		},
		cli.StringFlag{
			Name:   FLAG_TOKEN,
			EnvVar: ENVVAR_TOKEN,
		},
		cli.BoolFlag{
			Name:   FLAG_SKIP_VERIFY_SSL,
			EnvVar: ENVVAR_SKIP_VERIFY_SSL,
		},
		cli.BoolFlag{
			Name:   FLAG_SKIP_VERIFY_VERSION,
			EnvVar: ENVVAR_SKIP_VERIFY_VERSION,
		},
	}
}

type CLIConfig interface {
	Output() string
	Timeout() time.Duration
	Debug() bool
	Endpoint() string
	Token() string
	VerifySSL() bool
	VerifyVersion() bool
}

type ContextCLIConfig struct {
	C *cli.Context
}

func NewContextCLIConfig(c *cli.Context) *ContextCLIConfig {
	return &ContextCLIConfig{
		C: c,
	}
}

func (c *ContextCLIConfig) Validate() error {
	requiredVars := []string{
		FLAG_ENDPOINT,
		FLAG_TOKEN,
	}

	for _, name := range requiredVars {
		if !c.C.IsSet(name) {
			return fmt.Errorf("Required Variable %s is not set!", name)
		}
	}

	return nil
}

func (c *ContextCLIConfig) Output() string {
	return c.C.String(FLAG_OUTPUT)
}

func (c *ContextCLIConfig) Timeout() time.Duration {
	return c.C.Duration(FLAG_TIMEOUT)
}

func (c *ContextCLIConfig) Debug() bool {
	return c.C.Bool(FLAG_DEBUG)
}

func (c *ContextCLIConfig) Endpoint() string {
	return c.C.String(FLAG_ENDPOINT)
}

func (c *ContextCLIConfig) Token() string {
	return c.C.String(FLAG_TOKEN)
}

func (c *ContextCLIConfig) VerifySSL() bool {
	return !c.C.Bool(FLAG_SKIP_VERIFY_SSL)
}

func (c *ContextCLIConfig) VerifyVersion() bool {
	return !c.C.Bool(FLAG_SKIP_VERIFY_VERSION)
}
