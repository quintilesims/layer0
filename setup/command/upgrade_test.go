package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/quintilesims/layer0/setup/instance/mock_instance"
)

func TestUpgrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)
		mockInstance.EXPECT().
			Upgrade("v1.0.0", false).
			Return(nil)

		return mockInstance
	}

	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Upgrade(), "l0-setup upgrade name v1.0.0"); err != nil {
		t.Fatal(err)
	}
}

func TestUpgradeForce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)
		mockInstance.EXPECT().
			Upgrade("v1.0.0", true).
			Return(nil)

		return mockInstance
	}

	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Upgrade(), "l0-setup upgrade --force name v1.0.0"); err != nil {
		t.Fatal(err)
	}
}
