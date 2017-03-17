package decorators

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/quintilesims/layer0/common/waitutils"
	"strings"
	"time"
)

type Retry struct {
	Clock waitutils.Clock
}

func (this *Retry) shouldRetry(err error) bool {
	// As we discover more errors that indicate throttling, they should be added here
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			code := awsErr.Code()

			if code == "Throttling" || code == "ThrottlingException" {
				return true
			}

			message := strings.ToLower(awsErr.Message())
			if code == "ClientException" && strings.Contains(message, "too many concurrent attempts") {
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
		Delay:   5 * time.Second,
		Check:   check,
		Clock:   this.Clock,
	}

	return callObject.Wait()
}
