package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/quintilesims/layer0/setup/instance/mock_instance"
)

func TestSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)
		mockInstance.EXPECT().
			Set(map[string]interface{}{
				"k1": "v1",
				"k2": "v2",
			}).
			Return(nil)

		return mockInstance
	}

	input := "l0-setup set "
	input += "--input k1=v1 "
	input += "--input k2=v2 "
	input += "name"

	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Set(), input); err != nil {
		t.Fatal(err)
	}
}
