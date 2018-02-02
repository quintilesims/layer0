package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/quintilesims/layer0/setup/instance/mock_instance"
)

func TestDestroy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)
		mockInstance.EXPECT().
			Destroy(false).
			Return(nil)

		return mockInstance
	}

	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Destroy(), "l0-setup destroy name"); err != nil {
		t.Fatal(err)
	}
}

func TestDestroyForce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)
		mockInstance.EXPECT().
			Destroy(true).
			Return(nil)

		return mockInstance
	}

	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Destroy(), "l0-setup destroy --force name"); err != nil {
		t.Fatal(err)
	}
}
