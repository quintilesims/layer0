package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestDebugAdmin(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	base.Client.EXPECT().
		ReadConfig().
		Return(&models.APIConfig{}, nil)

	adminCommand := NewAdminCommand(base.Command())
	c := testutils.NewTestContext(t, nil, nil)

	if err := adminCommand.debug(c); err != nil {
		t.Fatal(err)
	}
}
