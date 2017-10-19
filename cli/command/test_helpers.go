package command

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/cli/printer"
	"github.com/quintilesims/layer0/cli/resolver/mock_resolver"
	"github.com/quintilesims/layer0/client/mock_client"
)

type TestCommandBase struct {
	Client   *mock_client.MockClient
	Printer  *printer.TestPrinter
	Resolver *mock_resolver.MockResolver
}

func NewTestCommand(t *testing.T) (*TestCommandBase, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	tc := &TestCommandBase{
		Client:   mock_client.NewMockClient(ctrl),
		Printer:  &printer.TestPrinter{},
		Resolver: mock_resolver.NewMockResolver(ctrl),
	}

	return tc, ctrl
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

func testWaitHelper(t *testing.T, fn func(t *testing.T, wait bool)) {
	t.Run("wait", func(t *testing.T) { fn(t, true) })
	t.Run("no-wait", func(t *testing.T) { fn(t, false) })
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
