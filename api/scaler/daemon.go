package scaler

import (
	"log"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/common/models"
)

func NewDaemonFN(jobStore job.Store, environmentProvider provider.EnvironmentProvider) func() error {
	return func() error {
		environments, err := environmentProvider.List()
		if err != nil {
			return err
		}

		for _, e := range environments {
			log.Printf("[DEBUG] [ScalerDaemon] Creating scale job for environment %s", e.EnvironmentID)
			if _, err := jobStore.Insert(models.ScaleEnvironmentJob, e.EnvironmentID); err != nil {
				return err
			}
		}

		return nil
	}
}
