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
			if job.Created.Before(time.Now().Add(-expiry)) {
				jobStore.Delete(job.JobID)
			}
		}

		return nil
	})
}
