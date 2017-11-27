package job

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type DynamoStore struct {
	table      dynamo.Table
	insertHook func(jobID string)
}

func NewDynamoStore(session *session.Session, table string) *DynamoStore {
	db := dynamo.New(session)

	return &DynamoStore{
		table:      db.Table(table),
		insertHook: func(string) {},
	}
}

func (d *DynamoStore) Clear() error {
	var jobs []models.Job
	if err := d.table.Scan().All(&jobs); err != nil {
		return err
	}

	for _, job := range jobs {
		if err := d.Delete(job.JobID); err != nil {
			return err
		}
	}

	return nil
}

func (d *DynamoStore) Insert(jobType models.JobType, req string) (string, error) {
	job := models.Job{
		JobID:   fmt.Sprintf("%v", time.Now().UnixNano()),
		Type:    jobType,
		Request: req,
		Status:  models.PendingJobStatus,
		Created: time.Now(),
		Result:  "",
	}

	if err := d.table.Put(job).Run(); err != nil {
		return "", err
	}

	d.insertHook(job.JobID)
	return job.JobID, nil
}

func (d *DynamoStore) SetInsertHook(hook func(jobID string)) {
	d.insertHook = hook
}

func (d *DynamoStore) AcquireJob(jobID string) (bool, error) {
	if err := d.table.Update("JobID", jobID).
		Set("Status", models.InProgressJobStatus).
		If("'Status' = ?", models.PendingJobStatus).
		Run(); err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (d *DynamoStore) Delete(jobID string) error {
	return d.table.Delete("JobID", jobID).Run()
}

func (d *DynamoStore) SetJobStatus(jobID string, status models.JobStatus) error {
	if err := d.table.Update("JobID", jobID).
		Set("Status", status).
		Run(); err != nil {
		return err
	}

	return nil
}

func (d *DynamoStore) SetJobResult(jobID, result string) error {
	if err := d.table.Update("JobID", jobID).
		Set("Result", result).
		Run(); err != nil {
		return err
	}

	return nil
}

func (d *DynamoStore) SetJobError(jobID string, err error) error {
	if err := d.table.Update("JobID", jobID).
		Set("Error", err.Error()).
		Set("Status", models.ErrorJobStatus).
		Run(); err != nil {
		return err
	}

	return nil
}

func (d *DynamoStore) SelectAll() ([]*models.Job, error) {
	jobs := []*models.Job{}
	if err := d.table.Scan().
		Consistent(false).
		All(&jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (d *DynamoStore) SelectByID(jobID string) (*models.Job, error) {
	var job *models.Job

	if err := d.table.Get("JobID", jobID).
		Consistent(true).
		One(&job); err != nil {

		if err.Error() == "dynamo: no item found" {
			return nil, errors.Newf(errors.JobDoesNotExist, "Job %s does not exist", jobID)
		}

		return nil, err
	}

	return job, nil
}
