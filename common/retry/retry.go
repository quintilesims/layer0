package retry

type RetryFunc func() (shouldRetry bool)

func Retry(fn RetryFunc, options ...Option) {
	for {
		for _, option := range options {
			if !option() {
				return
			}
		}
		shouldRetry := fn()
		if !shouldRetry {
			break
		}
	}

	return
}
