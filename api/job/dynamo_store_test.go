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
	table := os.Getenv(config.FlagTestAWSJobTable.EnvVar)
	if table == "" {
		t.Skipf("Test table not set (envvar: %s)", config.FlagTestAWSJobTable.EnvVar)
	}

	store := NewDynamoStore(session, table)
	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	return store
}

func TestDynamoStoreInsert(t *testing.T) {
	store := newTestStore(t)

	if _, err := store.Insert(models.DeleteEnvironmentJob, "1"); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoStoreInsertHook(t *testing.T) {
	store := newTestStore(t)

	var called bool
	store.SetInsertHook(func(jobID string) {
		called = true
	})

	if _, err := store.Insert(models.DeleteEnvironmentJob, "1"); err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
}

func TestDynamoStoreDelete(t *testing.T) {
	store := newTestStore(t)

	jobID, err := store.Insert(models.DeleteEnvironmentJob, "1")
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
		if _, err := store.Insert(models.DeleteEnvironmentJob, "1"); err != nil {
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
		{Type: models.DeleteEnvironmentJob, Request: "0"},
		{Type: models.DeleteEnvironmentJob, Request: "1"},
		{Type: models.DeleteServiceJob, Request: "2"},
		{Type: models.DeleteLoadBalancerJob, Request: "3"},
		{Type: models.DeleteTaskJob, Request: "4"},
	}

	for _, job := range jobs {
		jobID, err := store.Insert(models.JobType(job.Type), job.Request)
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

	jobID, err := store.Insert(models.DeleteEnvironmentJob, "1")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.SetJobStatus(jobID, models.ErrorJobStatus); err != nil {
		t.Fatal(err)
	}

	job, err := store.SelectByID(jobID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, models.ErrorJobStatus, job.Status)
}

func TestDynamoStoreSetResult(t *testing.T) {
	store := newTestStore(t)

	jobID, err := store.Insert(models.DeleteEnvironmentJob, "1")
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

	jobID, err := store.Insert(models.DeleteEnvironmentJob, "1")
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
