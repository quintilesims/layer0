package command

import (
	"testing"

	"github.com/golang/mock/gomock"
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

	commandFactory := NewCommandFactory(instanceFactory, nil)
	action := extractAction(t, commandFactory.Set())

	flags := map[string]interface{}{
		"input": []string{"k1=v1", "k2=v2"},
	}

	c := NewContext(t, []string{"name"}, flags)
	if err := action(c); err != nil {
		t.Fatal(err)
	}
}
