package testutils

/*
	There is a mock packge, https://github.com/DATA-DOG/go-sqlmock, that could be used for these tests
	instead of a real mysql database. However, testing against a real mysql database better suits
	our needs as we focus more heavily on integration testing rather than unit testing.
*/

import (
	"flag"
	"fmt"
	"github.com/quintilesims/layer0/common/db"
)

var (
	username = flag.String("username", "layer0", "username for the test db")
	password = flag.String("password", "nohaxplz", "password for the test db")
	address  = flag.String("address", "127.0.0.1", "address for the test db")
	port     = flag.Int("port", 3306, "port for the test db")
	dbName   = flag.String("dbname", "layer0_test", "name of the test db")
)

func init() {
	flag.Parse()
}

func GetDBConfig() db.Config {
	return db.Config{
		Connection: fmt.Sprintf("%s:%s@tcp(%s:%d)/", *username, *password, *address, *port),
		DBName:     *dbName,
	}
}
