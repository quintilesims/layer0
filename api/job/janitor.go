package job

import (
	"time"
)

func NewJanitorFN(jobStore Store, expiry time.Duration) func() error {
	return func() error {
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
	}
}
