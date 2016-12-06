package logic

import (
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/data"
	"github.com/quintilesims/layer0/common/models"
)

type Logic struct {
	Backend     backend.Backend
	AdminData   data.SQLAdmin
	TagData     data.TagData
	JobData     data.JobData
	ebStackName *string
}

func NewLogic(
	sqlAdmin data.SQLAdmin,
	tagData data.TagData,
	jobData data.JobData,
	backend backend.Backend,
) *Logic {
	return &Logic{
		AdminData: sqlAdmin,
		TagData:   tagData,
		JobData:   jobData,
		Backend:   backend,
	}
}

func (this *Logic) upsertTagf(entityID, entityType, key, value string) error {
	tag := models.EntityTag{
		EntityID:   entityID,
		EntityType: entityType,
		Key:        key,
		Value:      value,
	}

	return this.upsertTag(tag)
}

func (this *Logic) upsertTag(tag models.EntityTag) error {
	filter := map[string]string{
		"id":    tag.EntityID,
		"type":  tag.EntityType,
		"key":   tag.Key,
		"value": tag.Value,
	}

	entityTags, err := this.TagData.GetTags(filter)
	if err != nil {
		return err
	}

	if len(entityTags) == 0 {
		return this.TagData.Make(tag)
	}

	return nil
}

func (this *Logic) addTag(entityID, entityType, key, value string) error {
	tag := models.EntityTag{
		EntityID:   entityID,
		EntityType: entityType,
		Key:        key,
		Value:      value,
	}

	return this.TagData.Make(tag)
}

func (this *Logic) deleteEntityTags(entityID, entityType string) error {
	filter := map[string]string{
		"type": entityType,
		"id":   entityID,
	}

	tags, err := this.TagData.GetTags(filter)
	if err != nil {
		return err
	}

	for _, tag := range rangeTags(tags) {
		if err := this.TagData.Delete(tag); err != nil {
			return err
		}
	}

	return nil
}
