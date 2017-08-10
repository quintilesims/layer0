package lock

import "fmt"

type Lock interface {
	Acquire() error
	Release() error
}

type LockError struct {
	error
}

func LockIsAcquiredError() *LockError {
	return &LockError{
		fmt.Errorf("Lock is already acquired"),
	}
}

func IsAcquiredError(err error) bool {
	return err == LockIsAcquiredError()
}
