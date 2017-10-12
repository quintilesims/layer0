package command

import (
	"flag"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

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
	flagSet.String(config.FLAG_TIMEOUT, "15m", "")
	flagSet.Parse(args)
	return cli.NewContext(nil, flagSet, nil)
}
