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
	queue := make(chan models.Job)
	dispatcher := NewDispatcher(store, queue)

	for i := 0; i < 10; i++ {
		jobID, err := store.Insert(DeleteEnvironmentJob, strconv.Itoa(i))
		if err != nil {
			t.Fatal(err)
		}

		if i%3 == 0 {
			if err := store.SetJobStatus(jobID, Pending); err != nil {
				t.Fatal(err)
			}
		}
	}

	go func() {
		for {
			job := <-queue
			assert.Equal(t, Pending, Status(job.Status))
		}
	}()

	if err := dispatcher.Run(); err != nil {
		t.Fatal(err)
	}
}
