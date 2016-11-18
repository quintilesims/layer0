package data

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gitlab.imshealth.com/xfra/layer0/common/config"
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

// DataStore for Tag data in mysql

func NewTagMySQLDataStore() (*TagDataStoreMySQL, error) {
	lazy, err := NewLazySQL(
		func() (*sql.DB, error) {
			return initTagMySQL()
		},
	)
	if err != nil {
		return nil, err
	}
	return &TagDataStoreMySQL{
		lazyDb: lazy,
	}, nil
}

type TagDataStoreMySQL struct {
	lazyDb *lazySQL
}

func initTagMySQL() (*sql.DB, error) {
	cnxn := config.MySQLConnection()
	db, err := sql.Open("mysql", cnxn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (this *TagDataStoreMySQL) Close() {
	this.lazyDb.Close()
}

func (this *TagDataStoreMySQL) Select() ([]models.EntityTag, error) {
	return this.queryTags("")
}

func (this *TagDataStoreMySQL) SelectByType(entityType string) ([]models.EntityTag, error) {
	return this.queryTags("where entity_type=?", entityType)
}

func (this *TagDataStoreMySQL) SelectByTag(key, value string) ([]models.EntityTag, error) {
	return this.queryTags("where tag_key=? and tag_value=?",
		key, value)
}

func (this *TagDataStoreMySQL) SelectByTagPrefix(key, prefix string) ([]models.EntityTag, error) {
	prefix = this.escapeWilds(prefix)
	return this.queryTags("where tag_key=? and tag_value like ? ",
		key, prefix+"%")
}

func (this *TagDataStoreMySQL) SelectByTagKey(key string) ([]models.EntityTag, error) {
	return this.queryTags("where tag_key=? ", key)
}

func (this *TagDataStoreMySQL) SelectById(id string) ([]models.EntityTag, error) {
	return this.queryTags("where entity_id=?", id)
}

func (this *TagDataStoreMySQL) SelectByIdPrefix(idprefix string) ([]models.EntityTag, error) {
	prefix := this.escapeWilds(idprefix)
	return this.queryTags("where entity_id like ? ", prefix+"%")
}

func (this *TagDataStoreMySQL) Insert(tag models.EntityTag) error {
	return this.execTags("insert into tags (tag_key, tag_value, entity_type, entity_id) VALUES (?, ?, ?, ?)",
		tag.Key, tag.Value, tag.EntityType, tag.EntityID)
}

func (this *TagDataStoreMySQL) Delete(tag models.EntityTag) error {
	return this.execTags("delete from tags where tag_key=? and tag_value=? and entity_type=? and entity_id=?",
		tag.Key, tag.Value, tag.EntityType, tag.EntityID)
}

func (this *TagDataStoreMySQL) queryTags(where string, args ...interface{}) ([]models.EntityTag, error) {
	db, err := this.lazyDb.GetDB()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("select %s from tags ", allFields)
	stmt, err := db.Prepare(query + where)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return this.readEntityTags(rows)
}

func (this *TagDataStoreMySQL) execTags(query string, args ...interface{}) error {
	db, err := this.lazyDb.GetDB()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	// the _ argument is a resultset, which can be checked for the number of rows deleted
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}

func (this *TagDataStoreMySQL) escapeWilds(input string) string {
	input = strings.Replace(input, "\\", "\\\\", -1)
	input = strings.Replace(input, "%", "\\%", -1)
	input = strings.Replace(input, "_", "\\_", -1)
	return input
}

// rather than select *, give the fields an order consistent with readEntityTags
const allFields = "id, tag_key, tag_value, entity_type, entity_id"

func (this *TagDataStoreMySQL) readEntityTags(rows *sql.Rows) ([]models.EntityTag, error) {
	result := make([]models.EntityTag, 0, 2)
	for rows.Next() {
		var model = models.EntityTag{}
		var id int
		err := rows.Scan(&id,
			&model.Key,
			&model.Value,
			&model.EntityType,
			&model.EntityID)
		if err != nil {
			return nil, err
		}
		result = append(result, model)
	}
	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}
