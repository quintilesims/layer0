package data

import (
	"fmt"
	"regexp"

	"gitlab.imshealth.com/xfra/layer0/common/config"
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

type sqlVersion struct {
	Version string
	Message []string
}

type SQLAdmin interface {
	GetStatus() (*models.SQLVersion, error)
	UpdateSQL() error
}

type L0AdminLayer struct {
	DataStore AdminDataStore
	DbName    string
	UserName  string
	UserPass  string
}

func NewSQLAdminLayer(dataStore AdminDataStore) *L0AdminLayer {
	userName, userPass, dbName, _ := extractUser(config.MySQLConnection())
	return &L0AdminLayer{
		DataStore: dataStore,
		DbName:    dbName,
		UserName:  userName,
		UserPass:  userPass,
	}
}

func extractUser(cnxn string) (string, string, string, error) {
	userexp, err := regexp.Compile("^([^:]+):([^@]+).*/(.*)$")
	if err != nil {
		return "", "", "", err
	}
	strings := userexp.FindStringSubmatch(cnxn)
	if len(strings) < 3 {
		return "", "", "", fmt.Errorf("Failed to match username and password from user cnxn")
	}
	return strings[1], strings[2], strings[3], nil
}

func (this *L0AdminLayer) GetStatus() (*models.SQLVersion, error) {
	// logic:
	// try the master connection (ping)
	// -> no connection? 0.0
	// show the user, database
	// -> error 0.0, not found? 0.0
	// show the tags table
	// -> error / not found? 0.0
	// success? 0.1
	lines, err := this.DataStore.DescribeTables(this.DbName)
	if err != nil {
		return &models.SQLVersion{
			Version: "0.0",
			Message: []string{
				"Failed to connect the database, try updating sql (POST {'Version'='latest'} /sql)",
				err.Error(),
			},
		}, nil
	}

	return &models.SQLVersion{
		Version: "0.1",
		Message: lines,
	}, nil
}

func (this *L0AdminLayer) UpdateSQL() error {
	err := this.DataStore.CreateDatabase(this.DbName)
	if err != nil {
		return err
	}

	err = this.DataStore.CreateL0User(this.DbName, this.UserName, this.UserPass)
	if err != nil {
		return err
	}

	err = this.DataStore.CreateTagTable(this.DbName)
	if err != nil {
		return err
	}

	err = this.DataStore.CreateJobTable(this.DbName)
	if err != nil {
		return err
	}

	return nil
}
