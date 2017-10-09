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
			Usage:  "todo",
		},
		cli.StringFlag{
			Name:   FLAG_TOKEN,
			EnvVar: ENVVAR_TOKEN,
			Usage:  "todo",
		},
		cli.BoolFlag{
			Name:   FLAG_SKIP_VERIFY_SSL,
			EnvVar: ENVVAR_SKIP_VERIFY_SSL,
			Usage:  "todo",
		},
		cli.BoolFlag{
			Name:   FLAG_SKIP_VERIFY_VERSION,
			EnvVar: ENVVAR_SKIP_VERIFY_VERSION,
			Usage:  "todo",
		},
		cli.BoolFlag{
			Name:   FLAG_NO_WAIT,
			EnvVar: ENVVAR_NO_WAIT,
			Usage:  "todo",
		},
	}
}

func ValidateCLIContext(c *cli.Context) error {
	requiredVars := []string{
		FLAG_ENDPOINT,
		FLAG_TOKEN,
	}

	for _, name := range requiredVars {
		if !c.IsSet(name) {
			return fmt.Errorf("Required Variable %s is not set!", name)
		}
	}

	return nil
}
