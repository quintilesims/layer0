package tag

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
)

func TestJanitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagStore := NewMemoryStore()
	taskProvider := mock_provider.NewMockTaskProvider(ctrl)

	// todo: put some old tags in the store, mock the ListTasks call

	janitor := NewJanitor(tagStore, taskProvider)
	if err := janitor.Run(); err != nil {
		t.Fatal(err)
	}

	// todo: assert only old tags got deleted
}
