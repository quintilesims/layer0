package config

import (
	"github.com/urfave/cli"
)

func SetupFlags() []cli.Flag {
	return []cli.Flag{
		FlagDebug,
	}
}

func ValidateSetupContext(c *cli.Context) error {
	requiredFlags := []cli.Flag{}
	return ValidateContext(c, requiredFlags)
}
