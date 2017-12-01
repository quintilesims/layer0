package command

import (
	"flag"
	"strconv"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

type Args []string
type Flags map[string]interface{}
type Option func(t *testing.T, c *cli.Context, flagSet *flag.FlagSet)

func NewContext(t *testing.T, args Args, flags Flags, options ...Option) *cli.Context {
	flagSet := &flag.FlagSet{}
	c := cli.NewContext(&cli.App{}, flagSet, nil)

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
			t.Fatalf("Unexpected flag type for '%s'", key)
		}
	}

	// add default global flags
	flagSet.String(config.FLAG_OUTPUT, "text", "")
	flagSet.String(config.FLAG_TIMEOUT, "15m", "")
	flagSet.Bool(config.FLAG_NO_WAIT, false, "")

	for _, option := range options {
		option(t, c, flagSet)
	}

	if err := flagSet.Parse(args); err != nil {
		t.Fatal(err)
	}

	return c
}

func SetGlobalFlag(key, val string) Option {
	return func(t *testing.T, c *cli.Context, flagSet *flag.FlagSet) {
		if err := c.GlobalSet(key, val); err != nil {
			t.Fatal(err)
		}
	}
}

func SetNoWait(b bool) Option {
	return SetGlobalFlag(config.FLAG_NO_WAIT, strconv.FormatBool(b))
}

func SetTimeout(d time.Duration) Option {
	return SetGlobalFlag(config.FLAG_TIMEOUT, d.String())
}

func SetVersion(v string) Option {
	return func(t *testing.T, c *cli.Context, flagSet *flag.FlagSet) {
		c.App.Version = v
	}
}
