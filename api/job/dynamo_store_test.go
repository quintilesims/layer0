package job

/*
import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/quintilesims/layer0/common/models"
)

func NewTestJobStore(t *testing.T) *DynamoJobStore {
	table := config.TestDynamoJobTableName()
	if table == "" {
		t.Skipf("Skipping test: %s not set", config.TEST_AWS_JOB_DYNAMO_TABLE)
	}

	creds := credentials.NewStaticCredentials(config.AWSAccessKey(), config.AWSSecretKey(), "")
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AWSRegion()),
	}

	session := session.New(awsConfig)
	store := NewDynamoJobStore(session, table)

	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	return store
}

func TestDynamoJobStoreInsert(t *testing.T) {
	store := NewTestJobStore(t)

	job := &models.Job{JobID: "1", JobType: string(DeleteEnvironmentJob)}
	if err := store.Insert(job); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoJobStoreDelete(t *testing.T) {
	store := NewTestJobStore(t)

	job := &models.Job{JobID: "1", JobType: string(DeleteEnvironmentJob)}
	if err := store.Insert(job); err != nil {
		t.Fatal(err)
	}

	if err := store.Delete(job.JobID); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoJobStoreSelectAll(t *testing.T) {
	store := NewTestJobStore(t)

	jobs := []*models.Job{
		{JobID: "1", JobType: string(DeleteEnvironmentJob)},
		{JobID: "2", JobType: string(DeleteEnvironmentJob)},
		{JobID: "3", JobType: string(DeleteServiceJob)},
		{JobID: "4", JobType: string(DeleteLoadBalancerJob)},
		{JobID: "5", JobType: string(DeleteTaskJob)},
	}

	for _, job := range jobs {
		if err := store.Insert(job); err != nil {
			t.Fatal(err)
		}
	}

	result, err := store.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	if r, e := len(result), len(jobs); r != e {
		t.Fatalf("Result had %d jobs, expected %d", r, e)
	}
}

func TestDynamoJobStoreSelectByID(t *testing.T) {
	store := NewTestJobStore(t)

	jobs := []*models.Job{
		{JobID: "1", JobType: string(DeleteEnvironmentJob)},
		{JobID: "2", JobType: string(DeleteEnvironmentJob)},
		{JobID: "3", JobType: string(DeleteServiceJob)},
		{JobID: "4", JobType: string(DeleteLoadBalancerJob)},
		{JobID: "5", JobType: string(DeleteTaskJob)},
	}

	for _, job := range jobs {
		if err := store.Insert(job); err != nil {
			t.Fatal(err)
		}
	}

	result, err := store.SelectByID(jobs[2].JobID)
	if err != nil {
		t.Fatal(err)
	}

	if r, e := result.JobID, jobs[2].JobID; r != e {
		t.Fatalf("Result was %#v, expected %#v", r, e)
	}
}

func TestDynamoJobStoreUpdateStatus(t *testing.T) {
	store := NewTestJobStore(t)

	job := &models.Job{JobID: "1", JobStatus: string(Pending)}
	if err := store.Insert(job); err != nil {
		t.Fatal(err)
	}

	if err := store.UpdateJobStatus(job.JobID, InProgress); err != nil {
		t.Fatal(err)
	}

	result, err := store.SelectByID(job.JobID)
	if err != nil {
		t.Fatal(err)
	}

	if r, e := JobStatus(result.JobStatus), InProgress; r != e {
		t.Fatalf("Status was '%s', expected '%s'", r, e)
	}
}

func TestDynamoJobStoreSetMeta(t *testing.T) {
	store := NewTestJobStore(t)

	job := &models.Job{JobID: "1", Meta: map[string]string{"alpha": "1"}}
	if err := store.Insert(job); err != nil {
		t.Fatal(err)
	}

	meta := map[string]string{"beta": "2"}
	if err := store.SetJobMeta(job.JobID, meta); err != nil {
		t.Fatal(err)
	}

	result, err := store.SelectByID(job.JobID)
	if err != nil {
		t.Fatal(err)
	}

	if r, e := result.Meta, meta; !reflect.DeepEqual(r, e) {
		t.Fatalf("Status was '%s', expected '%s'", r, e)
	}

}
*/
