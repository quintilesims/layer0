package db

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

const CREATE_JOB_TABLE_QUERY = `CREATE TABLE IF NOT EXISTS jobs (
  job_id varchar(64) PRIMARY KEY NOT NULL,
  task_id varchar(64) NOT NULL,
  job_type smallint NOT NULL,
  job_status smallint NOT NULL,
  time_created timestamp DEFAULT '0000-00-00 00:00:00',
  last_updated timestamp NOT NULL ON UPDATE NOW(),
  request TEXT,
  meta TEXT
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
