package tag

import (
	"github.com/quintilesims/layer0/api/janitor"
	"github.com/quintilesims/layer0/api/provider"
)

func NewJanitor(tagStore Store, taskProvider provider.TaskProvider) *janitor.Janitor {
	return janitor.NewJanitor("Tag", func() error {
		// todo: delete all tags that don't exist in taskProvider
		return nil
	})
}
