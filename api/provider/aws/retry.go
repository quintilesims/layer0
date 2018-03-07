package aws

import (
	"time"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
)

func Retry(id string, describeFN func(id string) error) error {
	timeout := time.After(config.DefaultRetryTimeOut)
	tick := time.Tick(config.DefaultRetryWaitTime)

	for {
		select {
		case <-timeout:
			return errors.New(errors.TimeOut, nil)
		case <-tick:
			if err := describeFN(id); err == nil {
				return nil
			}
		}
	}
}
