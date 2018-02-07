package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestDebugAdmin(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewAdminCommand(base.CommandBase()).Command()

	base.Client.EXPECT().
		ReadConfig().
		Return(&models.APIConfig{}, nil)

	input := "l0 admin debug"
	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}
