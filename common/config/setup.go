package config

import (
	"fmt"

	"github.com/urfave/cli"
)

func SetupFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:   fmt.Sprintf("d, %s", FLAG_DEBUG),
			EnvVar: ENVVAR_DEBUG,
		},
	}
}
