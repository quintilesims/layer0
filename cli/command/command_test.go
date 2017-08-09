package command

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/client/mock_client"
	"github.com/quintilesims/layer0/cli/command/mock_command"
	"github.com/quintilesims/layer0/cli/printer"
)

type TestCommand struct {
	Client   *mock_client.MockClient
	Printer  *printer.TestPrinter
	Resolver *mock_command.MockResolver
}

func newTestCommand(t *testing.T) (*TestCommand, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	tc := &TestCommand{
		Client:   mock_client.NewMockClient(ctrl),
		Printer:  &printer.TestPrinter{},
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
