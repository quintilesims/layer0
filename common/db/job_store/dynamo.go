package job_store

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
)

type DynamoJobStore struct {
	table dynamo.Table
}

func NewDynamoJobStore(session *session.Session, table string) *DynamoJobStore {
	db := dynamo.New(session)

	return &DynamoJobStore{
		table: db.Table(table),
	}
}

func (d *DynamoJobStore) Init() error {
	return nil
}

func (d *DynamoJobStore) Clear() error {
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

func (d *DynamoJobStore) Insert(job *models.Job) error {
	return d.table.Put(job).Run()
}

func (d *DynamoJobStore) Delete(jobID string) error {
	return d.table.Delete("JobID", jobID).Run()
}

func (d *DynamoJobStore) UpdateJobStatus(jobID string, status types.JobStatus) error {

	if err := d.table.Update("JobID", jobID).Set("JobStatus", int64(status)).Run(); err != nil {
		return err
	}
	if status == types.Completed {
		//Add TTL for job entry
		if err := d.table.Update("JobID", jobID).Set("TimeToExist", time.Now().Add(time.Hour*time.Duration(config.DELETE_COMPLETED_JOB_TTL)).Unix()).Run(); err != nil {
			return err
		}
		job, err := d.SelectByID(jobID)
		if err != nil {
			return err
		}
		taskID := job.TaskID
		creds := credentials.NewStaticCredentials(config.AWSAccessKey(), config.AWSSecretKey(), "")
		session := session.New(config.GetAWSConfig(creds, config.AWSRegion()))
		tagStore := tag_store.NewDynamoTagStore(session, config.DynamoTagTableName())
		//Set TTL for task
		tagStore.AddTTLValue("task", taskID, config.TASK_TAG_TTL)
		//Set TTL for deploy starting with "job."
		tagStore.AddTTLValueToDeployJobs("task", taskID)
		//Set TTL for job
		tagStore.AddTTLValue("job", jobID, config.JOB_TAG_TTL)
	}
	return nil
}

func (d *DynamoJobStore) SetJobMeta(jobID string, meta map[string]string) error {
	if err := d.table.Update("JobID", jobID).Set("Meta", meta).Run(); err != nil {
		return err
	}

	return nil
}

func (d *DynamoJobStore) SelectAll() ([]*models.Job, error) {
	jobs := []*models.Job{}
	if err := d.table.Scan().
		Consistent(false).
		All(&jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (d *DynamoJobStore) SelectByID(jobID string) (*models.Job, error) {
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
