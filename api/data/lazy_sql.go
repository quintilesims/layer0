package data

import (
	"database/sql"
	"fmt"
)

// A SQL connection will fail to connect immediately if the credentials are
// invalid or if the target database does not exist.
// Because we want to create the database, tables, and users at runtime
// this class will retry the connection every time GetDB is called, and
// cache a successful result
type lazySQL struct {
	Db      *sql.DB
	Connect func() (*sql.DB, error)
}

func NewLazySQL(connectFunction func() (*sql.DB, error)) (*lazySQL, error) {
	if connectFunction == nil {
		return nil, fmt.Errorf("Connection Function was nil")
	}

	lazy := lazySQL{
		Connect: connectFunction,
	}

	return &lazy, nil
}

func (this *lazySQL) Close() {
	if this.Db != nil {
		this.Db.Close()
	}
}

func (this *lazySQL) GetDB() (*sql.DB, error) {
	err := this.refresh()
	if err != nil {
		return nil, err
	}
	return this.Db, nil
}

func (this *lazySQL) refresh() error {
	if this.Db == nil {
		db, err := this.Connect()
		if err != nil {
			return err
		}
		this.Db = db
	}

	err := this.Db.Ping()
	if err != nil {
		this.Db = nil
		return err
	}

	return nil
}
