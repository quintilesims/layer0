package job

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

var DeleteServiceSteps = []Step{
	Step{
		Name:    "Delete Service",
		Timeout: time.Minute * 10,
		Action:  DeleteService,
	},
}

func DeleteService(quit chan bool, context *JobContext) error {
	serviceID := context.Request()

	return runAndRetry(quit, time.Second*10, func() error {
		log.Infof("Running Action: DeleteService on '%s'", serviceID)
		return context.ServiceLogic.DeleteService(serviceID)
	})
}
