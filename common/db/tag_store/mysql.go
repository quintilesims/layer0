package tag_store

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/quintilesims/layer0/common/db/common"
	"github.com/quintilesims/layer0/common/models"
)

type MysqlTagStore struct {
	db     *sql.DB
	config common.Config
}

func NewMysqlTagStore(c common.Config) *MysqlTagStore {
	return &MysqlTagStore{
		config: c,
	}
}

// Creates the database and tag table if it doesn't already exist
func (m *MysqlTagStore) Init() error {
	db, err := sql.Open("mysql", m.config.Connection())
	if err != nil {
		return err
	}
	m.db = db

	if err := common.CreateDatabase(m.config.DBName, m.db); err != nil {
		return err
	}

	return m.exec(common.CREATE_TAG_TABLE_QUERY)
}

func (m *MysqlTagStore) Close() {
	m.db.Close()
}

func (m *MysqlTagStore) Clear() error {
	return m.exec("DELETE FROM tags")
}

func (m *MysqlTagStore) Insert(tag *models.Tag) error {
	return m.exec("INSERT INTO tags (tag_key, tag_value, entity_type, entity_id) VALUES (?, ?, ?, ?)",
		tag.Key,
		tag.Value,
		tag.EntityType,
		tag.EntityID)
}

func (m *MysqlTagStore) Delete(tag *models.Tag) error {
	return m.exec("DELETE FROM tags WHERE tag_key=? AND tag_value=? AND entity_type=? AND entity_id=?",
		tag.Key,
		tag.Value,
		tag.EntityType,
		tag.EntityID)
}

func (m *MysqlTagStore) SelectAll() (models.Tags, error) {
	return m.query("SELECT * FROM tags")
}

func (m *MysqlTagStore) SelectByEntityID(id string) (models.Tags, error) {
	return m.query("SELECT * FROM tags WHERE entity_id=?", id)
}

func (m *MysqlTagStore) SelectByEntityType(entityType string) (models.Tags, error) {
	return m.query("SELECT * FROM tags WHERE entity_type=?", entityType)
}

func (m *MysqlTagStore) exec(query string, args ...interface{}) error {
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(args...); err != nil {
		return err
	}

	return nil
}

func (m *MysqlTagStore) query(query string, args ...interface{}) (models.Tags, error) {
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := models.Tags{}
	for rows.Next() {
		var tag = models.Tag{}
		var id int

		err := rows.Scan(
			&id,
			&tag.Key,
			&tag.Value,
			&tag.EntityID,
			&tag.EntityType)
		if err != nil {
			return nil, err
		}

		tags = append(tags, &tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return models.Tags(tags), nil
}
