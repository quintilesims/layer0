package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestAdminDebug(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewAdminCommand(tc.Command())

	tc.Client.EXPECT().
		GetVersion().
		Return("v1.2.3", nil)

	c := testutils.GetCLIContext(t, nil, nil)
	if err := command.Debug(c); err != nil {
		t.Fatal(err)
	}
}

func TestAdminVersion(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewAdminCommand(tc.Command())

	tc.Client.EXPECT().
		GetVersion().
		Return("v1.2.3", nil)

	c := testutils.GetCLIContext(t, nil, nil)
	if err := command.Version(c); err != nil {
		t.Fatal(err)
	}
}

func TestAdminSQL(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewAdminCommand(tc.Command())

	tc.Client.EXPECT().
		UpdateSQL().
		Return(nil)

	c := testutils.GetCLIContext(t, nil, nil)
	if err := command.SQL(c); err != nil {
		t.Fatal(err)
	}
}

func TestAdminScale(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewAdminCommand(tc.Command())

	tc.Client.EXPECT().
		RunScaler("id").
		Return(&models.ScalerRunInfo{}, nil)

	tc.Resolver.EXPECT().
		Resolve("environment", "env").
		Return([]string{"id"}, nil)

	c := testutils.GetCLIContext(t, []string{"env"}, nil)
	if err := command.Scale(c); err != nil {
		t.Fatal(err)
	}
}
