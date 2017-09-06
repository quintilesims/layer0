package tftest

import (
	"fmt"
	"io"
)

type Logger interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

type IOLogger struct {
	writer io.Writer
}

func NewIOLogger(w io.Writer) *IOLogger {
	return &IOLogger{writer: w}
}

func (i *IOLogger) Log(args ...interface{}) {
	fmt.Fprint(i.writer, args...)
}

func (i *IOLogger) Logf(format string, args ...interface{}) {
	fmt.Fprintf(i.writer, format, args...)
}
