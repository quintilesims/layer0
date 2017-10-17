package command

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/cli/printer"
	"github.com/quintilesims/layer0/cli/resolver/mock_resolver"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

type TestCommandBase struct {
	Client   *mock_client.MockClient
	Printer  *printer.TestPrinter
	Resolver *mock_resolver.MockResolver
}

func newTestCommand(t *testing.T) (*TestCommandBase, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	tc := &TestCommandBase{
		Client:   mock_client.NewMockClient(ctrl),
		Printer:  &printer.TestPrinter{},
		Resolver: mock_resolver.NewMockResolver(ctrl),
	}

	return tc, ctrl
}

func (c *TestCommandBase) Command() *CommandBase {
	return &CommandBase{
		client:   c.Client,
		printer:  c.Printer,
		resolver: c.Resolver,
	}
}

type Args []string
type Flags map[string]interface{}

func getCLIContext(t *testing.T, args Args, flags Flags) *cli.Context {
	flagSet := &flag.FlagSet{}
	app := &cli.App{}

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
	flagSet.String(config.FLAG_OUTPUT, "text", "")
	flagSet.String(config.FLAG_TIMEOUT, "15m", "")
	flagSet.Parse(args)
	return cli.NewContext(app, flagSet, nil)
}

func createTempFile(t *testing.T, content string) (*os.File, func()) {
	file, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	return file, func() { os.Remove(file.Name()) }
}
