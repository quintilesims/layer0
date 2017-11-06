package job

import (
	"testing"
	"time"
)

func TestJanitor(t *testing.T) {
	store := NewMemoryStore()
	janitor := NewJanitor(store, time.Hour*24)

	// todo: put some jobs in the store

	if err := janitor.Run(); err != nil {
		t.Fatal(err)
	}

	// todo: assert only old jobs got deleted
}
