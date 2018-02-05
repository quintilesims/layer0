package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/quintilesims/layer0/setup/instance/mock_instance"
)

func TestPlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)
		mockInstance.EXPECT().
			Plan().
			Return(nil)

		return mockInstance
	}

	input := "l0-setup plan name"
	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Plan(), input); err != nil {
		t.Fatal(err)
	}
}
