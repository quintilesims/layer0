package job

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/lock/mock_lock"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestWorkerRunsJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLock := mock_lock.NewMockLock(ctrl)

	var called bool
	runner := RunnerFunc(func(j models.Job) (string, error) {
		called = true
		return "", nil
	})

	store := NewMemoryStore()
	jobID, err := store.Insert(models.DeleteEnvironmentJob, "1")
	if err != nil {
		t.Fatal(err)
	}

	mockLock.EXPECT().
		Acquire(jobID).
		Return(true, nil)

	mockLock.EXPECT().
		Release(jobID).
		Return(nil)

	queue := make(chan string)
	worker := NewWorker(0, store, queue, runner, mockLock)

	quit := worker.Start()
	defer quit()

	queue <- jobID
	assert.True(t, called)
}
