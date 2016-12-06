package decorators

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/quintilesims/layer0/common/waitutils"
	"time"
)

type Retry struct {
	Clock waitutils.Clock
}

func (this *Retry) shouldRetry(err error) bool {
	// As we discover more errors that indicate throttling, they should be added here
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// a simple case from this error report
			// https://github.com/ansible/ansible-modules-core/issues/143
			if awsErr.Code() == "Throttling" || awsErr.Code() == "ThrottlingException" {
				return true
			}
		}
	}

	return false
}

func (this *Retry) CallWithRetries(name string, call func() error) error {
	check := func() (bool, error) {
		err := call()
		if err == nil {
			return true, nil
		} else if this.shouldRetry(err) {
			return false, nil
		}

		return false, err
	}

	callObject := &waitutils.Waiter{
		Name:    name,
		Retries: 20,
		Delay:   10 * time.Second,
		Check:   check,
		Clock:   this.Clock,
	}

	return callObject.Wait()
}
