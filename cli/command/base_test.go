package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/cli/printer"
	"github.com/quintilesims/layer0/cli/resolver/mock_resolver"
	"github.com/quintilesims/layer0/client/mock_client"
)

type TestCommandBase struct {
	client   *mock_client.MockClient
	printer  *printer.TestPrinter
	resolver *mock_resolver.MockResolver
}

func newTestCommand(t *testing.T) (*TestCommandBase, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	tc := &TestCommandBase{
		client:   mock_client.NewMockClient(ctrl),
		printer:  &printer.TestPrinter{},
		resolver: mock_resolver.NewMockResolver(ctrl),
	}

	return tc, ctrl
}

func (tc *TestCommandBase) Command() *CommandBase {
	return &CommandBase{
		client:   tc.client,
		printer:  tc.printer,
		resolver: tc.resolver,
	}
}
