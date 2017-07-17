package decorators

import (
	"time"

	"github.com/quintilesims/layer0/common/logutils"
)

var log = logutils.NewStandardLogger("AWS Decorator")

func CallWithLogging(name string, call func() error) error {
	log.Debugf("AWS `%s` start", name)

	startTime := time.Now()
	err := call()
	duration := time.Since(startTime)

	if err != nil {
		log.Debugf("AWS `%s` Error: %v after %v", name, err, duration)
	}

	return err
}
