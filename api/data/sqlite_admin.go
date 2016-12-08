// +build !scratch

package data

import (
	"database/sql"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"github.com/quintilesims/layer0/common/config"
)

// NewSQLiteAdminDataStore is a DataStore for Admin data in sqlite
func NewSQLiteAdminDataStore() (*AdminDataStoreSQLite, error) {
	dbPath := config.SQLiteDbPath()
	if dbPath == "" {
		file, err := ioutil.TempFile("", "lite")
		if err != nil {
			return nil, err
		}
		dbPath = file.Name()
	}

	db, err := initAdminSQLite(dbPath)
	if err != nil {
		return nil, err
	}

	return &AdminDataStoreSQLite{
		Db:   db,
		File: dbPath,
	}, nil
}

func initAdminSQLite(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

type AdminDataStoreSQLite struct {
	Db   *sql.DB
	File string
}

func (this *AdminDataStoreSQLite) Close() {
	this.Db.Close()
	os.Remove(this.File)
}

// admin functions
func (this *AdminDataStoreSQLite) DescribeTables(dbName string) ([]string, error) {
	rows, err := this.Db.Query("SELECT sql FROM sqlite_master")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	log.Infof("Columns: %v", cols)

	result := make([]string, 0, 2)
	for rows.Next() {
		var row string
		err := rows.Scan(&row)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (this *AdminDataStoreSQLite) CreateDatabase(dbName string) error {
	return nil
}

func (this *AdminDataStoreSQLite) CreateTagTable(dbName string) error {
	return nil
}

func (this *AdminDataStoreSQLite) CreateJobTable(dbName string) error {
	return nil
}

func (this *AdminDataStoreSQLite) CreateL0User(dbName, username, password string) error {
	return nil
}
