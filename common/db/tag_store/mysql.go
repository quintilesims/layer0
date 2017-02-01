package tag_store

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	dbcommon "github.com/quintilesims/layer0/common/db"
	"github.com/quintilesims/layer0/common/models"
)

// Since query by anything other than EntityID and EntityType requires a mixture of
// and/or logic (depending on the context) and because the Name field can be any string,
// any complex queries required outside of entityType or EntityID need to be handled by the caller
type Query struct {
	EntityID   string
	EntityType string
}

type MysqlTagStore struct {
	db     *sql.DB
	config dbcommon.Config
}

func NewMysqlTagStore(c dbcommon.Config) *MysqlTagStore {
	return &MysqlTagStore{
		config: c,
	}
}

// Creates the database and tag table if it doesn't already exist
func (m *MysqlTagStore) Init() error {
	db, err := sql.Open("mysql", m.config.Connection)
	if err != nil {
		return err
	}
	m.db = db

	if err := dbcommon.CreateDatabase(m.config.DBName, m.db); err != nil {
		return err
	}

	return m.exec(dbcommon.CREATE_TAG_TABLE_QUERY)
}

func (m *MysqlTagStore) connect() (*sql.DB, error) {
	connection := fmt.Sprintf("%s%s", m.config.Connection, m.config.DBName)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
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

func (m *MysqlTagStore) SelectByQuery(entityType, entityID string) (models.Tags, error) {
	if entityType == "" && entityID == "" {
		return m.SelectAll()
	}

	clauses := []string{}
	args := []interface{}{}

	if entityType != "" {
		clauses = append(clauses, "entity_type=?")
		args = append(args, entityType)
	}

	if entityID != "" {
		clauses = append(clauses, "entity_id=?")
		args = append(args, entityID)
	}

	query := "SELECT * FROM tags WHERE "
	for i, clause := range clauses {
		if i == 0 {
			query += clause
		} else {
			query += fmt.Sprintf(" AND %s", clause)
		}
	}

	return m.query(query, args...)
}

func (m *MysqlTagStore) exec(query string, args ...interface{}) error {
	db, err := m.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
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
	db, err := m.connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
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
