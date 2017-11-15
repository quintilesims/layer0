package job

import (
	"fmt"
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
			fmt.Println("HERE", time.Since(job.Created), expiry)
			fmt.Println(time.Since(job.Created) > expiry)
			if time.Since(job.Created) > expiry {
				if err := jobStore.Delete(job.JobID); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
