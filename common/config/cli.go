package config

import (
	"fmt"

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

	for _, flag := range requiredFlags {
		name := flag.GetName()
		if !c.IsSet(name) {
			return fmt.Errorf("Required Variable %s is not set!", name)
		}
	}

	return nil
}
