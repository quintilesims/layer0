package command

import (
	"testing"
)

func TestAdminDebug(t *testing.T) {
	tc, ctrl := newTestCommand(t)
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

func TestAdminVersion(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewAdminCommand(tc.Command())

	tc.Client.EXPECT().
		GetVersion().
		Return("v1.2.3", nil)

	c := getCLIContext(t, nil, nil)
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

	c := getCLIContext(t, nil, nil)
	if err := command.SQL(c); err != nil {
		t.Fatal(err)
	}
}
