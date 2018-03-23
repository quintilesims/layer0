package retry

import "time"

type RetryFunc func() (shouldRetry bool, err error)

func Retry(fn RetryFunc, options ...Option) error {
	for {
		for _, option := range options {
			if err := option(); err != nil {
				return err
			}
		}

		shouldRetry, err := fn()
		if err != nil {
			return err
		}

		if !shouldRetry {
			break
		}
	}

	return nil
}

func SimpleRetry(fn RetryFunc, maxAttempts int, delay time.Duration) error {
	return Retry(fn, WithMaxAttempts(maxAttempts), WithDelay(delay))
}
