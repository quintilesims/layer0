package command

import (
	"net/url"
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

func TestReadLayer0InstanceLogs(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewAdminCommand(base.CommandBase()).Command()

	query := url.Values{
		"tail":  []string{"100"},
		"start": []string{"start"},
		"end":   []string{"end"},
	}

	base.Client.EXPECT().
		ReadLayer0InstanceLogs(query).
		Return([]models.LogFile{}, nil)

	input := "l0 admin instancelogs "
	input += "--tail 100 "
	input += "--start start "
	input += "--end end"

	if err := testutils.RunApp(command, input); err != nil {
		t.Fatal(err)
	}
}
