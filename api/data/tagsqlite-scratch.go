// +build scratch

package data

import (
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

// DataStore for Tag data in sqlite

func NewTagSQLiteDataStore() (*TagDataStoreSQLite, error) {
	return &TagDataStoreSQLite{}, nil
}

type TagDataStoreSQLite struct {
}

func (this *TagDataStoreSQLite) Close() {
}

// admin functions
func (this *TagDataStoreSQLite) DescribeTables(dbName string) ([]string, error) {
	return []string{}, nil
}

func (this *TagDataStoreSQLite) CreateDatabase(dbName string) error {
	return nil
}

func (this *TagDataStoreSQLite) CreateTagTable(dbName string) error {
	return nil
}

func (this *TagDataStoreSQLite) CreateL0User(dbName, username, password string) error {
	return nil
}

// tag functions

func (this *TagDataStoreSQLite) Select() ([]models.EntityTag, error) {
	return []models.EntityTag{}, nil
}

func (this *TagDataStoreSQLite) SelectByType(entityType string) ([]models.EntityTag, error) {
	return []models.EntityTag{}, nil
}

func (this *TagDataStoreSQLite) SelectByTag(key, value string) ([]models.EntityTag, error) {
	return []models.EntityTag{}, nil
}

func (this *TagDataStoreSQLite) SelectByTagPrefix(key, prefix string) ([]models.EntityTag, error) {
	return []models.EntityTag{}, nil
}

func (this *TagDataStoreSQLite) SelectByTagKey(key string) ([]models.EntityTag, error) {
	return []models.EntityTag{}, nil
}

func (this *TagDataStoreSQLite) SelectById(id string) ([]models.EntityTag, error) {
	return []models.EntityTag{}, nil
}

func (this *TagDataStoreSQLite) SelectByIdPrefix(idprefix string) ([]models.EntityTag, error) {
	return []models.EntityTag{}, nil
}

func (this *TagDataStoreSQLite) Insert(tag models.EntityTag) error {
	return nil
}

func (this *TagDataStoreSQLite) Delete(tag models.EntityTag) error {
	return nil
}
