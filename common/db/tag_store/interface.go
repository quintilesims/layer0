package tag_store

import (
	"github.com/quintilesims/layer0/common/models"
)

type TagStore interface {
	Init() error
	Close()
	Delete(tag *models.Tag) error
	Insert(tag *models.Tag) error
	SelectAll() models.Tags
	SelectByEntityID(string) models.Tags
	SelectByEntityType(string) models.Tags
}
