package job

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestWorkerRunsJob(t *testing.T) {
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

	queue := make(chan string)
	worker := NewWorker(0, store, queue, runner)

	quit := worker.Start()
	defer quit()

	queue <- jobID
	assert.True(t, called)
}
