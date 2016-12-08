// +build !scratch

package data

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
	"io/ioutil"
	"os"
	"time"
)

// NewJobSQLiteDataStore is a DataStore for Job data in sqlite
func NewJobSQLiteDataStore() (*JobDataStoreSQLite, error) {
	dbPath := config.SQLiteDbPath()
	if dbPath == "" {
		file, err := ioutil.TempFile("", "lite")
		if err != nil {
			return nil, err
		}
		dbPath = file.Name()
	}

	db, err := initSQLite(dbPath)
	if err != nil {
		return nil, err
	}

	return &JobDataStoreSQLite{
		Db:   db,
		File: dbPath,
	}, nil
}

type JobDataStoreSQLite struct {
	Db   *sql.DB
	File string
}

func initSQLite(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	initJobs := `
CREATE TABLE IF NOT EXISTS jobs (
	job_id text primary key,
	task_id text,
	job_type integer,
	job_status integer,
	time_created text,
	last_updated text,
	request blob,
	meta blob
);`

	if _, err := db.Exec(initJobs); err != nil {
		return nil, err
	}

	return db, nil
}

func (this *JobDataStoreSQLite) Close() {
	this.Db.Close()
	os.Remove(this.File)
}

// job functions
func (this *JobDataStoreSQLite) Select() ([]models.Job, error) {
	return this.queryJobs("select * from jobs")
}

func (this *JobDataStoreSQLite) SelectByID(id string) (*models.Job, error) {
	model, err := this.queryJobs("select * from jobs where job_id=?", id)
	if err != nil {
		return nil, err
	}

	return &model[0], nil
}

func (this *JobDataStoreSQLite) Insert(job *models.Job) error {
	query := `
		INSERT INTO jobs (job_id, task_id, job_status, job_type, request, time_created, last_updated, meta)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	meta, err := mapToString(job.Meta)
	if err != nil {
		return err
	}

	if job.TimeCreated.IsZero() {
		job.TimeCreated = time.Now()
	}

	if job.LastUpdated.IsZero() {
		job.LastUpdated = time.Now()
	}

	return this.execJobs(
		query,
		job.JobID,
		job.TaskID,
		job.JobStatus,
		job.JobType,
		job.Request,
		job.TimeCreated.Format(TIME_FORMAT),
		job.LastUpdated.Format(TIME_FORMAT),
		meta)
}

func (this *JobDataStoreSQLite) Delete(id string) error {
	return this.execJobs("delete from jobs where job_id=?", id)
}

func (this *JobDataStoreSQLite) UpdateStatus(jobID string, status types.JobStatus) error {
	query := `
		UPDATE jobs
		SET job_status=?,last_updated=?
		WHERE job_id=?`

	return this.execJobs(query, int64(status), time.Now().Format(TIME_FORMAT), jobID)
}

func (this *JobDataStoreSQLite) GetMeta(jobID string) (map[string]string, error) {
	model, err := this.SelectByID(jobID)
	if err != nil {
		return nil, err
	}

	return model.Meta, nil
}

func (this *JobDataStoreSQLite) SetMeta(jobID, key, val string) error {
	model, err := this.SelectByID(jobID)
	if err != nil {
		return err
	}

	model.Meta[key] = val

	meta, err := mapToString(model.Meta)
	if err != nil {
		return err
	}

	query := `
                UPDATE jobs
                SET meta=?,last_updated=?
                WHERE job_id=?`

	return this.execJobs(query, meta, time.Now().Format(TIME_FORMAT), jobID)
}

func (this *JobDataStoreSQLite) queryJobs(query string, args ...interface{}) ([]models.Job, error) {
	stmt, err := this.Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readJobs(rows)
}

func (this *JobDataStoreSQLite) execJobs(query string, args ...interface{}) error {
	stmt, err := this.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// the _ argument is a resultset, which can be checked for the number of rows deleted
	if _, err := stmt.Exec(args...); err != nil {
		return err
	}
	return nil
}

func readJobs(rows *sql.Rows) ([]models.Job, error) {
	result := make([]models.Job, 0, 2)
	for rows.Next() {
		var model = models.Job{}

		var created string
		var updated string
		var meta string

		if err := rows.Scan(
			&model.JobID,
			&model.TaskID,
			&model.JobType,
			&model.JobStatus,
			&created,
			&updated,
			&model.Request,
			&meta,
		); err != nil {
			return nil, err
		}

		if t, err := time.Parse(TIME_FORMAT, created); err != nil {
			return nil, err
		} else {
			model.TimeCreated = t
		}

		if t, err := time.Parse(TIME_FORMAT, updated); err != nil {
			return nil, err
		} else {
			model.LastUpdated = t
		}

		if m, err := stringToMap(meta); err != nil {
			return nil, err
		} else {
			model.Meta = m
		}

		log.Debugf("Scan: %d Model: %v", model.JobID, model)
		result = append(result, model)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
