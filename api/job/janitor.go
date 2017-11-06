package job

import (
	"time"

	"github.com/quintilesims/layer0/api/janitor"
)

func NewJanitor(jobStore Store, expiry time.Duration) *janitor.Janitor {
	return janitor.NewJanitor("Job", func() error {
		// todo: select all jobs, delete those older than expiry
		return nil
	})
}
