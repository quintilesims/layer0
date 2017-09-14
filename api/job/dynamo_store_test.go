package job

import (
	"fmt"
	"os"
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func newTestStore(t *testing.T) *DynamoStore {
	session := config.GetTestAWSSession()
	table := os.Getenv(config.ENVVAR_TEST_AWS_DYNAMO_JOB_TABLE)
	if table == "" {
		t.Skipf("Test table not set (envvar: %s)", config.ENVVAR_TEST_AWS_DYNAMO_JOB_TABLE)
	}

	store := NewDynamoStore(session, table)
	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	return store
}

func TestDynamoStoreInsert(t *testing.T) {
	store := newTestStore(t)

	if _, err := store.Insert(DeleteEnvironmentJob, "1"); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoStoreInsertHook(t *testing.T) {
	store := newTestStore(t)

	var called bool
	store.SetInsertHook(func(jobID string) {
		called = true
	})

	if _, err := store.Insert(DeleteEnvironmentJob, "1"); err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
}

func TestAcquireJobSuccess(t *testing.T) {
	store := newTestStore(t)

	jobID, err := store.Insert(DeleteEnvironmentJob, "1")
	if err != nil {
		t.Fatal(err)
	}

	ok, err := store.AcquireJob(jobID)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, ok)
}

func TestAcquireJobFailure(t *testing.T) {
	store := newTestStore(t)

	jobID, err := store.Insert(DeleteEnvironmentJob, "1")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.SetJobStatus(jobID, InProgress); err != nil {
		t.Fatal(err)
	}

	ok, err := store.AcquireJob(jobID)
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, ok)
}

func TestDynamoStoreDelete(t *testing.T) {
	store := newTestStore(t)

	jobID, err := store.Insert(DeleteEnvironmentJob, "1")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Delete(jobID); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoStoreSelectAll(t *testing.T) {
	store := newTestStore(t)

	for i := 0; i < 5; i++ {
		if _, err := store.Insert(DeleteEnvironmentJob, "1"); err != nil {
			t.Fatal(err)
		}
	}

	jobs, err := store.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, jobs, 5)
}

func TestDynamoStoreSelectByID(t *testing.T) {
	store := newTestStore(t)

	jobs := []*models.Job{
		{Type: string(DeleteEnvironmentJob), Request: "0"},
		{Type: string(DeleteEnvironmentJob), Request: "1"},
		{Type: string(DeleteServiceJob), Request: "2"},
		{Type: string(DeleteLoadBalancerJob), Request: "3"},
		{Type: string(DeleteTaskJob), Request: "4"},
	}

	for _, job := range jobs {
		jobID, err := store.Insert(JobType(job.Type), job.Request)
		if err != nil {
			t.Fatal(err)
		}

		result, err := store.SelectByID(jobID)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, job.Request, result.Request)
	}
}

func TestDynamoStoreSetJobStatus(t *testing.T) {
	store := newTestStore(t)

	jobID, err := store.Insert(DeleteEnvironmentJob, "1")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.SetJobStatus(jobID, Error); err != nil {
		t.Fatal(err)
	}

	job, err := store.SelectByID(jobID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, Error, Status(job.Status))
}

func TestDynamoStoreSetResult(t *testing.T) {
	store := newTestStore(t)

	jobID, err := store.Insert(DeleteEnvironmentJob, "1")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.SetJobResult(jobID, "result"); err != nil {
		t.Fatal(err)
	}

	job, err := store.SelectByID(jobID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "result", job.Result)
}

func TestDynamoStoreSetJobError(t *testing.T) {
	store := newTestStore(t)

	jobID, err := store.Insert(DeleteEnvironmentJob, "1")
	if err != nil {
		t.Fatal(err)
	}

	testError := fmt.Errorf("some error")
	if err := store.SetJobError(jobID, testError); err != nil {
		t.Fatal(err)
	}

	job, err := store.SelectByID(jobID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testError.Error(), job.Error)
}
