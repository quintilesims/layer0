package job

import (
	"time"

	"github.com/quintilesims/layer0/api/janitor"
)

func NewJanitor(jobStore Store, expiry time.Duration) *janitor.Janitor {
	return janitor.NewJanitor("Job", func() error {
		jobs, err := jobStore.SelectAll()
		if err != nil {
			return err
		}

		for _, job := range jobs {
			if time.Since(job.Created) > expiry {
				if err := jobStore.Delete(job.JobID); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
