package command

import (
	"flag"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/urfave/cli"
	"github.com/quintilesims/layer0/cli/client/mock_client"
	"github.com/quintilesims/layer0/cli/command/mock_command"
	"github.com/quintilesims/layer0/cli/printer"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var TEST_TIMEOUT = time.Minute * 15

type Args []string

type Flags map[string]interface{}

func getCLIContext(t *testing.T, args []string, flags Flags) *cli.Context {
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
	flagSet.String("timeout", fmt.Sprintf("%v", TEST_TIMEOUT), "")
	flagSet.Parse(args)
	return cli.NewContext(nil, flagSet, nil)
}

type TestCommand struct {
	Client   *mock_client.MockClient
	Printer  *printer.FakePrinter
	Resolver *mock_command.MockResolver
}

func newTestCommand(t *testing.T) (*TestCommand, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	tc := &TestCommand{
		Client:   mock_client.NewMockClient(ctrl),
		Printer:  &printer.FakePrinter{},
		Resolver: mock_command.NewMockResolver(ctrl),
	}

	return tc, ctrl
}

func (tc *TestCommand) Command() *Command {
	return &Command{
		Client:   tc.Client,
		Printer:  tc.Printer,
		Resolver: tc.Resolver,
	}
}

// todo: place in testutils
func tempFile(t *testing.T, content string) (*os.File, func()) {
	file, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	return file, func() { os.Remove(file.Name()) }
}
