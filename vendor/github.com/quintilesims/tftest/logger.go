package tftest

import (
	"testing"
)

type Logger interface {
	Printf(string, ...interface{})
}

type TestLogger struct {
	t *testing.T
}

func NewTestLogger(t *testing.T) *TestLogger {
	return &TestLogger{t: t}
}

func (l *TestLogger) Printf(format string, tokens ...interface{}) {
	l.t.Logf(format, tokens...)
}
