package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/quintilesims/layer0/setup/instance/mock_instance"
)

func TestUnset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)
		mockInstance.EXPECT().
			Unset("k1").
			Return(nil)

		return mockInstance
	}

	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Unset(), "l0-setup unset name k1"); err != nil {
		t.Fatal(err)
	}
}
