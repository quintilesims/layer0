package common

import (
	"database/sql"
	"fmt"
)

const CREATE_TAG_TABLE_QUERY = `CREATE TABLE IF NOT EXISTS tags (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  tag_key varchar(64),
  tag_value varchar(64),
  entity_id varchar(64),
  entity_type varchar(32),
  INDEX ix_tag_value (tag_key(10), tag_value(10)),
  INDEX ix_id (entity_id(10))
);`

func CreateDatabase(name string, db *sql.DB) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", name)
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = fmt.Sprintf("USE %s", name)
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
