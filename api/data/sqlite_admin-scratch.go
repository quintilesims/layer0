// +build scratch

package data

import (
	"database/sql"
	"fmt"
	"os"
)

// NewSQLiteAdminDataStore is a DataStore for Admin data in sqlite
func NewSQLiteAdminDataStore() (*AdminDataStoreSQLite, error) {
	return nil, fmt.Errorf("Not implemented")
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
	return nil, nil
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
