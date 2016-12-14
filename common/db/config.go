package db

import (
	"fmt"
)

type Config struct {
	Username string
	Password string
	Address  string
	Port     int
	DBName   string
}

func (c Config) Connection() string {
	// format: 'username:password@tcp(127.0.0.1:3306)/[dbname?]'
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		c.Username,
		c.Password,
		c.Address,
		c.Port)
}
