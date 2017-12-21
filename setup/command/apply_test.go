package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/config"
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

	commandFactory := NewCommandFactory(instanceFactory, nil)
	action := extractAction(t, commandFactory.Apply())

	c := config.NewTestContext(t, []string{"name"}, nil)
	if err := action(c); err != nil {
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

	commandFactory := NewCommandFactory(instanceFactory, nil)
	action := extractAction(t, commandFactory.Apply())

	c := config.NewTestContext(t, []string{"name"}, map[string]interface{}{"quick": "true"})
	if err := action(c); err != nil {
		t.Fatal(err)
	}
}
