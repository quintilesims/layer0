package job

import (
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.FatalLevel)

	timeMultiplier = 0
}
