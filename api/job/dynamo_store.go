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
	table dynamo.Table
}

func NewDynamoStore(session *session.Session, table string) *DynamoStore {
	db := dynamo.New(session)

	return &DynamoStore{
		table: db.Table(table),
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

func (d *DynamoStore) Insert(jobType JobType, req string) (string, error) {
	job := models.Job{
		JobID:   fmt.Sprintf("%v", time.Now().UnixNano()),
		Type:    string(jobType),
		Request: req,
		Status:  string(Pending),
		Created: time.Now(),
		Meta:    map[string]string{},
	}

	if err := d.table.Put(job).Run(); err != nil {
		return "", err
	}

	return job.JobID, nil
}

func (d *DynamoStore) AcquireJob(jobID string) (bool, error) {
	if err := d.table.Update("JobID", jobID).
		Set("Status", InProgress).
		If("'Status' = ?", Pending).
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

func (d *DynamoStore) SetJobStatus(jobID string, status Status) error {
	if err := d.table.Update("JobID", jobID).
		Set("Status", status).
		Run(); err != nil {
		return err
	}

	return nil
}

func (d *DynamoStore) SetJobMeta(jobID string, meta map[string]string) error {
	if err := d.table.Update("JobID", jobID).
		Set("Meta", meta).
		Run(); err != nil {
		return err
	}

	return nil
}

func (d *DynamoStore) SetJobError(jobID string, err error) error {
	if err := d.table.Update("JobID", jobID).
		Set("Error", err.Error()).
		Set("Status", Error).
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
