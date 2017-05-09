package tag_store

import (
	"github.com/quintilesims/layer0/common/models"
)

type TagStore interface {
	Init() error
	Delete(entityType, entityID, key string) error
	Insert(tag models.Tag) error
	SelectByType(entityType string) (models.Tags, error)
	SelectByTypeAndID(entityType, entityID string) (models.Tags, error)
}
