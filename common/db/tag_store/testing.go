package tag_store

/*
	There is a mock packge, https://github.com/DATA-DOG/go-sqlmock, that could be used for these tests
	instead of a real mysql database. However, testing against a real mysql database better suits
	our needs as we focus more heavily on integration testing rather than unit testing.
*/

import (
	"flag"
	"github.com/quintilesims/layer0/common/db/common"
	"github.com/quintilesims/layer0/common/models"
	"testing"
)

// todo: may have to place this into common/testutils for shared use w/ job data
const (
	DEFAULT_USERNAME = "layer0"
	DEFAULT_PASSWORD = "nohaxplz"
	DEFAULT_ADDRESS  = "127.0.0.1"
	DEFAULT_PORT     = 3306
	DEFAULT_DB_NAME  = "layer0_test"
)

var (
	username = flag.String("username", DEFAULT_USERNAME, "username for the test db")
	password = flag.String("password", DEFAULT_PASSWORD, "password for the test db")
	address  = flag.String("address", DEFAULT_ADDRESS, "address for the test db")
	port     = flag.Int("port", DEFAULT_PORT, "port for the test db")
	dbName   = flag.String("dbname", DEFAULT_DB_NAME, "name of the test db")
)

func init() {
	flag.Parse()
}

func NewTestTagStore(t *testing.T) *MysqlTagStore {
	store := NewMysqlTagStore(common.Config{
		Username: *username,
		Password: *password,
		Address:  *address,
		Port:     *port,
		DBName:   *dbName,
	})

	if err := store.Init(); err != nil {
		t.Fatal(err)
	}

	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	return store
}

func NewTestTagStoreWithTags(t *testing.T, tags models.Tags) *MysqlTagStore {
	store := NewTestTagStore(t)
	for _, tag := range tags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	return store
}
