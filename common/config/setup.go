package config

import (
	"fmt"

	"github.com/urfave/cli"
)

func SetupFlags() []cli.Flag {
	return []cli.Flag{
		FlagDebug,
	}
}

func ValidateSetupContext(c *cli.Context) error {
	requiredFlags := []cli.Flag{}

	for _, flag := range requiredFlags {
		name := flag.GetName()
		if !c.IsSet(name) {
			return fmt.Errorf("Required Variable %s is not set!", name)
		}
	}

	return nil
}
