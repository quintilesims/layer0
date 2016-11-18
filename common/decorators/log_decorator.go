package decorators

import (
	"gitlab.imshealth.com/xfra/layer0/common/logutils"
	"time"
)

var log = logutils.NewStackTraceLogger("AWS Decorator")

func CallWithLogging(name string, call func() error) error {
	log.Debugf("AWS `%s` start", name)

	startTime := time.Now()
	err := call()
	duration := time.Since(startTime)

	if err != nil {
		log.Debugf("AWS `%s` Error: %v after %v", name, err, duration)
	} else {
		log.Debugf("AWS `%s` call complete after %v", name, duration)
	}

	return err
}
