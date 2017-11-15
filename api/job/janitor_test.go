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

	expiry := -24 * time.Hour
	now := time.Now()

	jobs := []*models.Job{
		{
			JobID:   "delete",
			Created: now.Add(expiry),
		},
		{
			JobID:   "keep",
			Created: now.Add(expiry * -2),
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
			Created: now.Add(expiry * -2),
		},
	}

	actual, err := store.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}
