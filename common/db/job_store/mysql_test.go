package job_store

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/common/types"
	"testing"
	"time"
)

func getTestJobs() []*models.Job {
	return []*models.Job{
		{
			JobID:       "job_id1",
			TaskID:      "task_id1",
			JobStatus:   int64(types.InProgress),
			JobType:     int64(types.DeleteEnvironmentJob),
			Request:     "request1",
			TimeCreated: time.Now(),
			LastUpdated: time.Now(),
			Meta:        map[string]string{"k1": "v1"},
		},
		{
			JobID:       "job_id2",
			TaskID:      "task_id2",
			JobStatus:   int64(types.Pending),
			JobType:     int64(types.DeleteServiceJob),
			Request:     "request2",
			TimeCreated: time.Now(),
			LastUpdated: time.Now(),
			Meta:        map[string]string{"k2": "v2"},
		},
		{
			JobID:       "job_id3",
			TaskID:      "task_id3",
			JobStatus:   int64(types.Completed),
			JobType:     int64(types.DeleteLoadBalancerJob),
			Request:     "request3",
			TimeCreated: time.Now(),
			LastUpdated: time.Now(),
			Meta:        map[string]string{"k3": "v3"},
		},
		{
			JobID:       "job_id4",
			TaskID:      "task_id4",
			JobStatus:   int64(types.Error),
			JobType:     int64(types.DeleteLoadBalancerJob),
			Request:     "request4",
			TimeCreated: time.Now(),
			LastUpdated: time.Now(),
			Meta:        map[string]string{"k4": "v4"},
		},
	}
}

func NewTestJobStore(t *testing.T) *MysqlJobStore {
	store := NewMysqlJobStore(testutils.GetDBConfig())

	if err := store.Init(); err != nil {
		t.Fatal(err)
	}

	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	return store
}

func NewTestJobStoreWithJobs(t *testing.T, jobs []*models.Job) *MysqlJobStore {
	store := NewTestJobStore(t)
	for _, job := range jobs {
		if err := store.Insert(job); err != nil {
			t.Fatal(err)
		}
	}

	return store
}

func assertJobsEqual(t *testing.T, target, expected *models.Job) {
	// don't compare times
	testutils.AssertEqual(t, target.JobID, expected.JobID)
	testutils.AssertEqual(t, target.TaskID, expected.TaskID)
	testutils.AssertEqual(t, target.JobStatus, expected.JobStatus)
	testutils.AssertEqual(t, target.JobType, expected.JobType)
	testutils.AssertEqual(t, target.Request, expected.Request)
	testutils.AssertEqual(t, target.Meta, expected.Meta)

}

func assertJobsMatch(t *testing.T, store *MysqlJobStore, expected []*models.Job) {
	jobs, err := store.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(jobs), len(expected))
	for i, job := range jobs {
		assertJobsEqual(t, job, expected[i])
	}
}

func TestMysqlJobStoreInsert(t *testing.T) {
	store := NewTestJobStore(t)
	defer store.Close()

	jobs := getTestJobs()
	for _, job := range jobs {
		if err := store.Insert(job); err != nil {
			t.Fatal(err)
		}
	}

	assertJobsMatch(t, store, jobs)
}

func TestMysqlJobStoreDelete(t *testing.T) {
	jobs := getTestJobs()
	store := NewTestJobStoreWithJobs(t, jobs)
	defer store.Close()

	for _, job := range jobs[:2] {
		if err := store.Delete(job.JobID); err != nil {
			t.Fatal(err)
		}
	}

	// calling delete on a non-existing id shouldn't thow an error
	if err := store.Delete("invalid"); err != nil {
		t.Fatal(err)
	}

	assertJobsMatch(t, store, jobs[2:])
}

func TestMysqlJobStoreUpdateJobStatus(t *testing.T) {
	jobs := getTestJobs()
	store := NewTestJobStoreWithJobs(t, jobs)
	defer store.Close()

	if err := store.UpdateJobStatus(jobs[0].JobID, types.Completed); err != nil {
		t.Fatal(err)
	}

	jobs[0].JobStatus = int64(types.Completed)
	assertJobsMatch(t, store, jobs)
}

func TestMysqlJobStoreSelectAll(t *testing.T) {
	jobs := getTestJobs()
	store := NewTestJobStoreWithJobs(t, jobs)
	defer store.Close()

	assertJobsMatch(t, store, jobs)
}

func TestMysqlJobStoreSelectByID(t *testing.T) {
	jobs := getTestJobs()
	store := NewTestJobStoreWithJobs(t, jobs)
	defer store.Close()

	for _, job := range jobs {
		j, err := store.SelectByID(job.JobID)
		if err != nil {
			t.Fatal(err)
		}

		assertJobsEqual(t, j, job)
	}
}

func TestMysqlJobStoreSelectByID_error(t *testing.T) {
	jobs := getTestJobs()
	store := NewTestJobStoreWithJobs(t, jobs)
	defer store.Close()

	if _, err := store.SelectByID("invalid"); err == nil {
		t.Fatalf("Error was nil!")
	}
}
