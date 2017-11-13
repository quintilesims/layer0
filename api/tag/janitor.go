package tag

import (
	"github.com/quintilesims/layer0/api/janitor"
	"github.com/quintilesims/layer0/api/provider"
)

func NewJanitor(tagStore Store, taskProvider provider.TaskProvider) *janitor.Janitor {
	return janitor.NewJanitor("Tag", func() error {
		tasks, err := taskProvider.List()
		if err != nil {
			return err
		}

		tags, err := tagStore.SelectByType("task")
		if err != nil {
			return err
		}

		m := make(map[string]bool, len(tasks))
		for _, task := range tasks {
			m[task.TaskID] = true
		}

		for _, tag := range tags {
			if !m[tag.EntityID] {
				tagStore.Delete(tag.EntityType, tag.EntityID, tag.Key)
			}
		}
		return nil
	})
}
