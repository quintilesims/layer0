package logic

import (
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/scheduler"
	"github.com/quintilesims/layer0/common/db/job_store"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
)

type Logic struct {
	Backend  backend.Backend
	TagStore tag_store.TagStore
	JobStore job_store.JobStore
	Scaler   scheduler.EnvironmentScaler
}

func NewLogic(
	tagStore tag_store.TagStore,
	jobData job_store.JobStore,
	backend backend.Backend,
	scaler scheduler.EnvironmentScaler,
) *Logic {
	return &Logic{
		TagStore: tagStore,
		JobStore: jobData,
		Backend:  backend,
		Scaler:   scaler,
	}
}

func (this *Logic) upsertTag(tag models.Tag) error {
	tags, err := this.TagStore.SelectByQuery(tag.EntityType, tag.EntityID)
	if err != nil {
		return err
	}

	exists := tags.Any(func(t models.Tag) bool {
		return t.Key == tag.Key && t.Value == tag.Value
	})

	if !exists {
		return this.TagStore.Insert(&tag)
	}

	return nil
}

func (this *Logic) deleteEntityTags(entityType, entityID string) error {
	tags, err := this.TagStore.SelectByQuery(entityType, entityID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := this.TagStore.Delete(tag.TagID); err != nil {
			return err
		}
	}

	return nil
}
