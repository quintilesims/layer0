package job

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestWorkerRunsJob(t *testing.T) {
	var called bool
	runner := func(j models.Job) error {
		called = true
		return nil
	}

	queue := make(chan models.Job)
	worker := NewWorker(0, queue, runner)

	quit := worker.Start()
	defer quit()

	queue <- models.Job{}
	assert.True(t, called)
}
