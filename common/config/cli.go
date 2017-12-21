package config

import (
	"github.com/urfave/cli"
)

func CLIFlags() []cli.Flag {
	return []cli.Flag{
		FlagDebug,
		FlagEndpoint,
		FlagToken,
		FlagOutput,
		FlagTimeout,
		FlagNoWait,
		FlagSkipVerifySSL,
		FlagSkipVerifyVersion,
	}
}

func ValidateCLIContext(c *cli.Context) error {
	requiredFlags := []cli.Flag{
		FlagToken,
	}

	return ValidateContext(c, requiredFlags)
}
