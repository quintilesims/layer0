package job

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDispatcherQueuesPendingJobs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMemoryStore()
	queue := make(chan string)
	dispatcher := NewDispatcher(store, queue)

	jobStatuses := map[string]models.JobStatus{}
	for i := 0; i < 50; i++ {
		jobID, err := store.Insert(models.DeleteEnvironmentJob, strconv.Itoa(i))
		if err != nil {
			t.Fatal(err)
		}

		switch i % 4 {
		case 0:
			jobStatuses[jobID] = models.Pending
		case 1:
			jobStatuses[jobID] = models.InProgress
		case 2:
			jobStatuses[jobID] = models.Completed
		case 3:
			jobStatuses[jobID] = models.Error

		}
	}

	for jobID, status := range jobStatuses {
		if err := store.SetJobStatus(jobID, status); err != nil {
			t.Fatal(err)
		}
	}

	go func() {
		for {
			jobID := <-queue
			assert.Equal(t, models.Pending, jobStatuses[jobID])
		}
	}()

	if err := dispatcher.Run(); err != nil {
		t.Fatal(err)
	}
}
