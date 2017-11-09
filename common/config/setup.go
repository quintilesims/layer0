package config

import "github.com/urfave/cli"

func SetupFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:   FLAG_DEBUG,
			EnvVar: ENVVAR_DEBUG,
		},
	}
}
