package lock

import "errors"

type Lock interface {
	Acquire() error
	Release() error
}

type ContentionError struct {
	error
}

func NewContentionError() *ContentionError {
	return &ContentionError{
		errors.New("Lock is in contention"),
	}
}

func IsContentionError(err error) bool {
	if _, ok := err.(*ContentionError); ok {
		return true
	}

	return false
}
