package tag

import (
	"github.com/quintilesims/layer0/common/models"
)

type Store interface {
	Delete(entityType, entityID, key string) error
	Insert(tag models.Tag) error
	SelectAll() (models.Tags, error)
	SelectByType(entityType string) (models.Tags, error)
	SelectByTypeAndID(entityType, entityID string) (models.Tags, error)
}
