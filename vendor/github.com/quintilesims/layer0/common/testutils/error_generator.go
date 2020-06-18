package testutils

import (
	"sync"
)

type ErrorGenerator struct {
	Index  int
	errors map[int]error
	once   sync.Once
}

func (e *ErrorGenerator) init() {
	e.errors = map[int]error{}
}

func (e *ErrorGenerator) Set(i int, err error) {
	e.once.Do(e.init)
	e.errors[i] = err
}

func (e *ErrorGenerator) Error() error {
	e.once.Do(e.init)
	e.Index++

	return e.errors[e.Index]
}
