package job_store

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
)

type DynamoJobStore struct {
	table          dynamo.Table
	consistentRead bool
}

func NewDynamoJobStore(session *session.Session, table string) *DynamoJobStore {
	db := dynamo.New(session)

	return &DynamoJobStore{
		table:          db.Table(table),
		consistentRead: false,
	}
}

func (d *DynamoJobStore) setConsistentRead(v bool) {
	d.consistentRead = v
}

func (d *DynamoJobStore) Init() error {
	d.setConsistentRead(true)
	return nil
}

func (d *DynamoJobStore) Clear() error {
	jobs, err := d.SelectAll()
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if err := d.Delete(job.JobID); err != nil {
			return err
		}
	}

	return nil
}

func (d *DynamoJobStore) Insert(job *models.Job) error {
	// ensure we perform a consistent read after a write
	defer d.setConsistentRead(true)
	return d.table.Put(job).Run()
}

func (d *DynamoJobStore) Delete(jobID string) error {
	// ensure we perform a consistent read after a delete
	defer d.setConsistentRead(true)
	return d.table.Delete("JobID", jobID).Run()
}

func (d *DynamoJobStore) UpdateJobStatus(jobID string, status types.JobStatus) error {
	// ensure we perform a consistent read after a write
	defer d.setConsistentRead(true)

	if err := d.table.Update("JobID", jobID).Set("JobStatus", int64(status)).Run(); err != nil {
		return err
	}

	return nil
}

func (d *DynamoJobStore) SetJobMeta(jobID string, meta map[string]string) error {
	// ensure we perform a consistent read after a write
	defer d.setConsistentRead(true)

	if err := d.table.Update("JobID", jobID).Set("Meta", meta).Run(); err != nil {
		return err
	}

	return nil
}

func (d *DynamoJobStore) SelectAll() ([]*models.Job, error) {
	defer d.setConsistentRead(false)

	jobs := []*models.Job{}
	if err := d.table.Scan().Consistent(true).All(&jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (d *DynamoJobStore) SelectByID(jobID string) (*models.Job, error) {
	defer d.setConsistentRead(false)

	var job *models.Job
	if err := d.table.Get("JobID", jobID).Consistent(true).One(&job); err != nil {
		return nil, err
	}

	return job, nil
}
