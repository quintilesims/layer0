package job

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDispatcherQueuesPendingJobs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMemoryStore()
	queue := make(chan models.Job)
	dispatcher := NewDispatcher(store, queue)

	jobs := []models.Job{
		{JobID: "j1", JobStatus: string(Pending)},
		{JobID: "j2", JobStatus: string(Pending)},
		{JobID: "j3", JobStatus: string(InProgress)},
		{JobID: "j4", JobStatus: string(Completed)},
	}

	for _, job := range jobs {
		store.Insert(job)
	}

	go func() {
		for {
			job := <-queue
			assert.Equal(t, Pending, Status(job.JobStatus))
		}
	}()

	if err := dispatcher.Run(); err != nil {
		t.Fatal(err)
	}
}
