// +build scratch

package data

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
)

// DataStore for Job data in sqlite

func NewJobSQLiteDataStore() (*JobDataStoreSQLite, error) {
	return &JobDataStoreSQLite{}, nil
}

type JobDataStoreSQLite struct {
}

func (this *JobDataStoreSQLite) Close() {
}

// admin functions
func (this *JobDataStoreSQLite) DescribeTables(dbName string) ([]string, error) {
	return []string{}, nil
}

// func (this *JobDataStoreSQLite) CreateDatabase(dbName string) error {
//	return nil
// }

func (this *JobDataStoreSQLite) CreateTagTable(dbName string) error {
	return nil
}

func (this *JobDataStoreSQLite) CreateL0User(dbName, username, password string) error {
	return nil
}

// job functions

func (this *JobDataStoreSQLite) Select() ([]models.Job, error) {
	return []models.Job{}, nil
}

func (this *JobDataStoreSQLite) SelectByID(jobID string) (*models.Job, error) {
	return &models.Job{}, nil
}

func (this *JobDataStoreSQLite) Insert(job *models.Job) error {
	return nil
}

func (this *JobDataStoreSQLite) Delete(jobID string) error {
	return nil
}

func (this *JobDataStoreSQLite) UpdateStatus(jobID string, status types.JobStatus) error {
	return nil
}

func (this *JobDataStoreSQLite) GetMeta(jobID string) (map[string]string, error) {
	return nil, nil
}

func (this *JobDataStoreSQLite) SetMeta(jobID, key, val string) error {
	return nil
}
