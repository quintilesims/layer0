// +build !scratch

package data

import (
	"database/sql"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.imshealth.com/xfra/layer0/common/config"
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

// NewTagSQLiteDataStore is a DataStore for Tag data in sqlite
func NewTagSQLiteDataStore() (*TagDataStoreSQLite, error) {
	dbPath := config.SQLiteDbPath()
	if dbPath == "" {
		file, err := ioutil.TempFile("", "lite")
		if err != nil {
			return nil, err
		}
		dbPath = file.Name()
	}

	db, err := initTagSQLite(dbPath)
	if err != nil {
		return nil, err
	}
	return &TagDataStoreSQLite{
		Db:   db,
		File: dbPath,
	}, nil
}

type TagDataStoreSQLite struct {
	Db   *sql.DB
	File string
}

func initTagSQLite(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	create table if not exists tags (
 		id integer primary key autoincrement,
 		tag_key text,
 		tag_value text,
 		entity_type text,
 		entity_id text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (this *TagDataStoreSQLite) Close() {
	this.Db.Close()
	os.Remove(this.File)
}

// tag functions
func (this *TagDataStoreSQLite) Select() ([]models.EntityTag, error) {
	return this.queryTags("select * from tags")
}

func (this *TagDataStoreSQLite) SelectByType(entityType string) ([]models.EntityTag, error) {
	return this.queryTags("select * from tags where entity_type=?", entityType)
}

func (this *TagDataStoreSQLite) SelectByTag(key, value string) ([]models.EntityTag, error) {
	return this.queryTags("select * from tags where tag_key=? and tag_value=?",
		key, value)
}

func (this *TagDataStoreSQLite) SelectByTagPrefix(key, prefix string) ([]models.EntityTag, error) {
	prefix = escapeWilds(prefix)
	return this.queryTags("select * from tags where tag_key=? and tag_value like ? ESCAPE '\\' ",
		key, prefix+"%")
}

func (this *TagDataStoreSQLite) SelectByTagKey(key string) ([]models.EntityTag, error) {
	return this.queryTags("select * from tags where tag_key=? ", key)
}

func (this *TagDataStoreSQLite) SelectById(id string) ([]models.EntityTag, error) {
	return this.queryTags("select * from tags where entity_id=?", id)
}

func (this *TagDataStoreSQLite) SelectByIdPrefix(idprefix string) ([]models.EntityTag, error) {
	prefix := escapeWilds(idprefix)
	return this.queryTags("select * from tags where entity_id like ? ESCAPE '\\' ",
		prefix+"%")
}

func escapeWilds(input string) string {
	input = strings.Replace(input, "\\", "\\\\", -1)
	input = strings.Replace(input, "%", "\\%", -1)
	input = strings.Replace(input, "_", "\\_", -1)
	return input
}

func (this *TagDataStoreSQLite) Insert(tag models.EntityTag) error {
	return this.execTags("insert into tags (tag_key, tag_value, entity_type, entity_id) VALUES  (?, ?, ?, ?)",
		tag.Key, tag.Value, tag.EntityType, tag.EntityID)
}

func (this *TagDataStoreSQLite) Delete(tag models.EntityTag) error {
	return this.execTags("delete from tags where tag_key=? and tag_value=? and entity_type=? and entity_id=?",
		tag.Key, tag.Value, tag.EntityType, tag.EntityID)
}

func (this *TagDataStoreSQLite) queryTags(query string, args ...interface{}) ([]models.EntityTag, error) {
	stmt, err := this.Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readEntityTags(rows)
}

func (this *TagDataStoreSQLite) execTags(query string, args ...interface{}) error {
	stmt, err := this.Db.Prepare(query)
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

func readEntityTags(rows *sql.Rows) ([]models.EntityTag, error) {
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
		log.Debugf("Scan: %d Model: %v", id, model)
		result = append(result, model)
	}
	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}
