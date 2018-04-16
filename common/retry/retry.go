package retry

type RetryFunc func() (shouldRetry bool)

func Retry(shouldRetryFN RetryFunc, options ...Option) error {
	for {
		for _, option := range options {
			if err := option(); err != nil {
				return err
			}
		}

		if !shouldRetryFN() {
			break
		}
	}

	return nil
}
