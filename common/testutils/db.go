package testutils

/*
	There is a mock packge, https://github.com/DATA-DOG/go-sqlmock, that could be used for these tests
	instead of a real mysql database. However, testing against a real mysql database better suits
	our needs as we focus more heavily on integration testing rather than unit testing.
*/

import (
	"flag"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db"
)

var (
	username = flag.String("username", config.DBUsername(), "username for the test db")
	password = flag.String("password", config.DBPassword(), "password for the test db")
	address  = flag.String("address", config.DBAddress(), "address for the test db")
	port     = flag.Int("port", config.DBPort(), "port for the test db")
	dbName   = flag.String("dbname", "layer0_test", "name of the test db")
)

func init() {
	flag.Parse()
}

func GetDBConfig() db.Config {
	return db.Config{
		Username: *username,
		Password: *password,
		Address:  *address,
		Port:     *port,
		DBName:   *dbName,
	}
}
