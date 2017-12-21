package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

func TestDebugAdmin(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Client.EXPECT().
		ReadConfig().
		Return(&models.APIConfig{}, nil)

	adminCommand := NewAdminCommand(base.Command())
	c := config.NewTestContext(t, nil, nil)

	if err := adminCommand.debug(c); err != nil {
		t.Fatal(err)
	}
}
