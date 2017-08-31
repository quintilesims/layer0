package tag

import (
	"github.com/quintilesims/layer0/common/models"
)

type Store interface {
	Init() error
	Delete(entityType, entityID, key string) error
	Insert(tag models.Tag) error
	SelectByType(entityType string) (models.Tags, error)
	SelectByTypeAndID(entityType, entityID string) (models.Tags, error)
}
