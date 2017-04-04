package job_store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	dbcommon "github.com/quintilesims/layer0/common/db"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
	"time"
)

// use pre-defined date/times to format https://golang.org/pkg/time/#Parse
const TIME_FORMAT = "2006-01-02 15:04:05"

type MysqlJobStore struct {
	config dbcommon.Config
}

func NewMysqlJobStore(c dbcommon.Config) *MysqlJobStore {
	return &MysqlJobStore{
		config: c,
	}
}

// Creates the database and job table if it doesn't already exist
func (m *MysqlJobStore) Init() error {
	db, err := sql.Open("mysql", m.config.Connection)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := dbcommon.CreateDatabase(m.config.DBName, db); err != nil {
		return err
	}

	if err := m.exec(dbcommon.CREATE_JOB_TABLE_QUERY); err != nil {
		return err
	}

	return nil
}

func (m *MysqlJobStore) connect() (*sql.DB, error) {
	connection := fmt.Sprintf("%s%s", m.config.Connection, m.config.DBName)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (m *MysqlJobStore) Clear() error {
	return m.exec("DELETE FROM jobs")
}

func (m *MysqlJobStore) Insert(job *models.Job) error {
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

	return m.exec(
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

func (m *MysqlJobStore) Delete(id string) error {
	return m.exec("DELETE FROM jobs WHERE job_id=?", id)
}

func (m *MysqlJobStore) UpdateJobStatus(id string, status types.JobStatus) error {
	return m.exec("UPDATE jobs SET job_status=? WHERE job_id=?", int64(status), id)
}

func (m *MysqlJobStore) SetJobMeta(id string, meta map[string]string) error {
	metaStr, err := mapToString(meta)
	if err != nil {
		return err
	}

	return m.exec("UPDATE jobs SET meta=? WHERE job_id=?", metaStr, id)
}

func (m *MysqlJobStore) SelectAll() ([]*models.Job, error) {
	return m.query("SELECT * FROM jobs")
}

func (m *MysqlJobStore) SelectByID(id string) (*models.Job, error) {
	jobs, err := m.query("SELECT * FROM jobs where job_id=?", id)
	if err != nil {
		return nil, err
	}

	if len(jobs) == 0 {
		return nil, fmt.Errorf("Job with id '%s' not found", id)
	}

	return jobs[0], nil
}

func (m *MysqlJobStore) exec(query string, args ...interface{}) error {
	db, err := m.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(args...); err != nil {
		return err
	}

	return nil
}

func (m *MysqlJobStore) query(query string, args ...interface{}) ([]*models.Job, error) {
	db, err := m.connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []*models.Job{}
	for rows.Next() {
		var job = models.Job{}
		var created string
		var updated string
		var meta string

		if err := rows.Scan(
			&job.JobID,
			&job.TaskID,
			&job.JobType,
			&job.JobStatus,
			&created,
			&updated,
			&job.Request,
			&meta,
		); err != nil {
			return nil, err
		}

		if t, err := time.Parse(TIME_FORMAT, created); err != nil {
			return nil, err
		} else {
			job.TimeCreated = t
		}

		if t, err := time.Parse(TIME_FORMAT, updated); err != nil {
			return nil, err
		} else {
			job.LastUpdated = t
		}

		if m, err := stringToMap(meta); err != nil {
			return nil, err
		} else {
			job.Meta = m
		}

		jobs = append(jobs, &job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
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
