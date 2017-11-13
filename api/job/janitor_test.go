package job

import (
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestJanitor(t *testing.T) {
	store := NewMemoryStore()
	janitor := NewJanitor(store, time.Hour*24)
	timeNow := time.Now()

	jobs := []*models.Job{
		{
			JobID:   "delete",
			Created: timeNow.Add(-24 * time.Hour),
		},
		{
			JobID:   "keep",
			Created: timeNow.Add(-12 * time.Hour),
		},
	}

	for _, job := range jobs {
		store.jobs = append(store.jobs, job)
	}

	if err := janitor.Run(); err != nil {
		t.Fatal(err)
	}

	expected := []*models.Job{
		{
			JobID:   "keep",
			Created: timeNow.Add(-12 * time.Hour),
		},
	}

	actual, _ := store.SelectAll()

	assert.Equal(t, expected, actual)
}
