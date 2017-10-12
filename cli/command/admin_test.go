package command

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestDebugAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)

	tc := &TestCommand{
		Client:   mock_client.NewMockClient(ctrl),
		Printer:  &printer.TestPrinter{},
		Resolver: mock_command.NewMockResolver(ctrl),
	}

	defer ctrl.Finish()
	command := NewAdminCommand(tc.Command())

	tc.Client.EXPECT().
		GetVersion().
		Return("v1.2.3", nil)

	c := getCLIContext(t, nil, nil)
	if err := command.Debug(c); err != nil {
		t.Fatal(err)
	}
}
