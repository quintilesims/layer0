package data

import (
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

type TagDataStore interface {
	Select() ([]models.EntityTag, error)
	SelectByType(entityType string) ([]models.EntityTag, error)
	SelectByTag(name, value string) ([]models.EntityTag, error)
	SelectByTagPrefix(name, prefix string) ([]models.EntityTag, error)
	SelectByTagKey(name string) ([]models.EntityTag, error)
	SelectById(id string) ([]models.EntityTag, error)
	SelectByIdPrefix(idprefix string) ([]models.EntityTag, error)
	Insert(tagDetail models.EntityTag) error
	Delete(tagDetail models.EntityTag) error
}
