package testutils

import (
	"flag"
	"github.com/urfave/cli"
	"testing"
	"time"
)

var TEST_TIMEOUT = time.Minute * 15

func GetCLIContext(t *testing.T, args []string, flags map[string]interface{}) *cli.Context {
	flagSet := &flag.FlagSet{}

	for key, val := range flags {
		switch v := val.(type) {
		case bool:
			flagSet.Bool(key, v, "")
		case string:
			flagSet.String(key, v, "")
		case []string:
			slice := cli.StringSlice(v)
			flagSet.Var(&slice, key, "")
		case int:
			flagSet.Int(key, v, "")
		default:
			t.Errorf("Cannot generate CLI context: unknown flag type for '%s'", key)
		}
	}

	// add default global flags
	flagSet.String("output", "text", "")
	flagSet.String("timeout", TEST_TIMEOUT.String(), "")
	flagSet.Parse(args)
	return cli.NewContext(nil, flagSet, nil)
}
