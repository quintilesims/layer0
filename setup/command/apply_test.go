package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/quintilesims/layer0/setup/instance/mock_instance"
)

func TestApply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)
		mockInstance.EXPECT().
			Apply(true).
			Return(nil)

		return mockInstance
	}

	input := "l0-setup apply "
	input += "--push=false "
	input += "name"

	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Apply(), input); err != nil {
		t.Fatal(err)
	}
}

func TestApplyQuick(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)
		mockInstance.EXPECT().
			Apply(false).
			Return(nil)

		return mockInstance
	}

	input := "l0-setup apply "
	input += "--push=false "
	input += "--quick "
	input += "name"

	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Apply(), input); err != nil {
		t.Fatal(err)
	}
}
