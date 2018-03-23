package retry

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
