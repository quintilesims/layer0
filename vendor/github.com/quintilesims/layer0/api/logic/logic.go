package logic

import (
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/scheduler"
	"github.com/quintilesims/layer0/common/db/job_store"
	"github.com/quintilesims/layer0/common/db/tag_store"
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

func (this *Logic) deleteEntityTags(entityType, entityID string) error {
	tags, err := this.TagStore.SelectByTypeAndID(entityType, entityID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := this.TagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
			return err
		}
	}

	return nil
}
