package data

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.imshealth.com/xfra/layer0/common/config"
)

type AdminDataStore interface {
	DescribeTables(dbName string) ([]string, error)
	CreateDatabase(dbName string) error
	CreateTagTable(dbName string) error
	CreateJobTable(dbName string) error
	CreateL0User(dbName, username, password string) error
}

func NewMySQLAdmin() (*L0MySQLAdmin, error) {
	connectMaster := func() (*sql.DB, error) {
		return initAdminMySQL("")
	}

	lazyMaster, err := NewLazySQL(connectMaster)
	if err != nil {
		return nil, err
	}

	return &L0MySQLAdmin{
		lazyMaster:   lazyMaster,
		lazySpecific: make(map[string]*lazySQL),
	}, nil
}

type L0MySQLAdmin struct {
	lazyMaster   *lazySQL
	lazySpecific map[string]*lazySQL
}

func initAdminMySQL(dbName string) (*sql.DB, error) {
	cnxn := config.MySQLAdminConnection()

	// replace the /{dbName} suffix with a custom dbname
	rx, err := regexp.Compile("/.*$")
	if err != nil {
		return nil, err
	}
	cnxn = rx.ReplaceAllString(cnxn, "/"+dbName)

	return sql.Open("mysql", cnxn)
}

func (this *L0MySQLAdmin) Close() {
	this.lazyMaster.Close()
	for _, v := range this.lazySpecific {
		v.Close()
	}
}

func (this *L0MySQLAdmin) getLazy(dbName string) (*lazySQL, error) {
	if _, ok := this.lazySpecific[dbName]; !ok {
		connect := func() (*sql.DB, error) {
			return initAdminMySQL(dbName)
		}

		lazy, err := NewLazySQL(connect)
		if err != nil {
			return nil, err
		}

		this.lazySpecific[dbName] = lazy
	}

	return this.lazySpecific[dbName], nil
}

func (this *L0MySQLAdmin) exec(lazyDb *lazySQL, query string, args ...interface{}) error {
	db, err := lazyDb.GetDB()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(query)
	log.Infof("Executing: %s", query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(args...); err != nil {
		return err
	}

	return nil
}

func (this *L0MySQLAdmin) DescribeTables(dbName string) ([]string, error) {
	lazy, err := this.getLazy(dbName)
	if err != nil {
		return nil, err
	}

	db, err := lazy.GetDB()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("explain tags")
	log.Infof("Describe tables query: [%s]", query)
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := []string{
		strings.Join(cols, " | "),
	}

	for rows.Next() {
		// http://stackoverflow.com/questions/14477941/read-select-columns-into-string-in-go
		readCols := make([]interface{}, len(cols))
		writeCols := make([]sql.NullString, len(cols))
		for i, _ := range writeCols {
			readCols[i] = &writeCols[i]
		}

		if err := rows.Scan(readCols...); err != nil {
			return nil, err
		}

		ss := nullStringToString(writeCols)
		stringRow := strings.Join(ss, " | ")
		result = append(result, stringRow)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func nullStringToString(nullResult []sql.NullString) []string {
	result := make([]string, len(nullResult))
	for i, v := range nullResult {
		if v.Valid {
			result[i] = v.String
		} else {
			result[i] = ""
		}
	}

	return result
}

func (this *L0MySQLAdmin) CreateDatabase(dbName string) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName)
	return this.exec(this.lazyMaster, query)
}

func (this *L0MySQLAdmin) CreateTagTable(dbName string) error {
	lazy, err := this.getLazy(dbName)
	if err != nil {
		return err
	}

	query := `
CREATE TABLE IF NOT EXISTS tags (
	id INTEGER AUTO_INCREMENT PRIMARY KEY,
	tag_key varchar(64),
	tag_value varchar(64),
	entity_id varchar(64),
	entity_type varchar(32),
	INDEX ix_tag_value (tag_key(10), tag_value(10)),
	INDEX ix_id (entity_id(10))
);`

	return this.exec(lazy, query)
}

func (this *L0MySQLAdmin) CreateJobTable(dbName string) error {
	lazy, err := this.getLazy(dbName)
	if err != nil {
		return err
	}

	query := `
CREATE TABLE IF NOT EXISTS jobs (
        job_id varchar(64) PRIMARY KEY NOT NULL,
        task_id varchar(64) NOT NULL,
        job_type smallint NOT NULL,
        job_status smallint NOT NULL,
        time_created timestamp DEFAULT '0000-00-00 00:00:00',
        last_updated timestamp NOT NULL ON UPDATE NOW(),
        request TEXT,
        meta TEXT
);`

	if err := this.exec(lazy, query); err != nil {
		return err
	}

	columns_071 := []string{
		"time_created timestamp DEFAULT '0000-00-00 00:00:00'",
		"last_updated timestamp NOT NULL ON UPDATE NOW()",
		"request TEXT",
		"meta TEXT",
	}

	for _, col := range columns_071 {
		query := fmt.Sprintf("ALTER TABLE jobs ADD %s", col)
		if err := this.exec(lazy, query); err != nil {
			if !strings.Contains(err.Error(), "Duplicate column name") {
				return err
			}
		}
	}

	return nil
}

func (this *L0MySQLAdmin) CreateL0User(dbName, username, password string) error {
	lazy, err := this.getLazy(dbName)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s';`, dbName, username, password)
	return this.exec(lazy, query)
}
