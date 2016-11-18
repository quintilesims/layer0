package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.imshealth.com/xfra/layer0/common/config"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/types"
	"time"
)

// NewJobMySQLDataStore is for Job data in mysql
func NewJobMySQLDataStore() (*JobDataStoreMySQL, error) {
	lazy, err := NewLazySQL(
		func() (*sql.DB, error) {
			return initJobMySQL()
		},
	)

	if err != nil {
		return nil, err
	}
	return &JobDataStoreMySQL{
		lazyDb: lazy,
	}, nil
}

type JobDataStoreMySQL struct {
	lazyDb *lazySQL
}

func initJobMySQL() (*sql.DB, error) {
	cnxn := config.MySQLConnection()
	db, err := sql.Open("mysql", cnxn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (this *JobDataStoreMySQL) Close() {
	this.lazyDb.Close()
}

func (this *JobDataStoreMySQL) Select() ([]models.Job, error) {
	return this.queryJobs("")
}

func (this *JobDataStoreMySQL) SelectByID(jobID string) (*models.Job, error) {
	models, err := this.queryJobs("where job_id=?", jobID)
	if err != nil {
		return nil, err
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("job with id '%s' not found", jobID)
	}

	return &models[0], nil
}

func (this *JobDataStoreMySQL) Insert(job *models.Job) error {
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

func (this *JobDataStoreMySQL) Delete(jobID string) error {
	if err := this.execJobs("delete from jobs where job_id=?", jobID); err != nil {
		return err
	}

	return nil
}

func (this *JobDataStoreMySQL) UpdateStatus(jobID string, jobStatus types.JobStatus) error {
	query := `
		UPDATE jobs
		SET job_status=?
		WHERE job_id=?`

	return this.execJobs(query, int64(jobStatus), jobID)
}

func (this *JobDataStoreMySQL) GetMeta(jobID string) (map[string]string, error) {
	model, err := this.SelectByID(jobID)
	if err != nil {
		return nil, err
	}

	return model.Meta, nil
}

func (this *JobDataStoreMySQL) SetMeta(jobID, key, val string) error {
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
                SET meta=?
                WHERE job_id=?`

	return this.execJobs(query, meta, jobID)
}

const allJobFields = "job_id, task_id, job_status, job_type, request, time_created, last_updated, meta"

func (this *JobDataStoreMySQL) queryJobs(where string, args ...interface{}) ([]models.Job, error) {
	db, err := this.lazyDb.GetDB()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("select %s from jobs ", allJobFields)
	stmt, err := db.Prepare(query + where)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return this.readJobs(rows)
}

func (this *JobDataStoreMySQL) execJobs(query string, args ...interface{}) error {
	db, err := this.lazyDb.GetDB()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(query)
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

// use pre-defined date/times to format https://golang.org/pkg/time/#Parse
const TIME_FORMAT = "2006-01-02 15:04:05"

func (this *JobDataStoreMySQL) readJobs(rows *sql.Rows) ([]models.Job, error) {
	result := make([]models.Job, 0, 2)

	for rows.Next() {
		var model = models.Job{}
		var created string
		var updated string
		var meta string

		if err := rows.Scan(
			&model.JobID,
			&model.TaskID,
			&model.JobStatus,
			&model.JobType,
			&model.Request,
			&created,
			&updated,
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

		result = append(result, model)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func mapToString(meta map[string]string) (string, error) {
	bytes, err := json.Marshal(meta)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func stringToMap(metaStr string) (map[string]string, error) {
	var meta map[string]string
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		return nil, err
	}

	if meta == nil {
		meta = map[string]string{}
	}

	return meta, nil
}
