package logutils

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/facebookgo/stack"
	"os"
	"strings"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

func SetGlobalLogger(logger *logrus.Logger) {
	logrus.SetFormatter(logger.Formatter)
}

type StandardFormatter struct {
	Name string
}

func NewStandardLogger(name string) *logrus.Logger {
	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: &StandardFormatter{Name: name},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.GetLevel(),
	}

	return logger
}

func (this *StandardFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	time := entry.Time.Format(TIME_FORMAT)
	level := strings.ToUpper(entry.Level.String())
	if level == "WARNING" {
		level = "WARN"
	}

	content := fmt.Sprintf("%v [%s] %-5s: %s\n", time, this.Name, level, entry.Message)
	return []byte(content), nil
}

type CLIFormatter struct {}

func (c *CLIFormatter) Format(entry *logrus.Entry) ([]byte, error) {
        level := strings.ToUpper(entry.Level.String())
        content := fmt.Sprintf("[%-5s] %s\n", level, entry.Message)
        return []byte(content), nil
}


type StackTraceFormatter struct {
	*StandardFormatter
}

func NewStackTraceLogger(name string) *logrus.Logger {
	logger := NewStandardLogger(name)
	logger.Formatter = &StackTraceFormatter{logger.Formatter.(*StandardFormatter)}

	return logger
}

func (this *StackTraceFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// logrus levels go backwards: Panic=0, Debug=5
	if entry.Level > logrus.ErrorLevel || logrus.GetLevel() != logrus.DebugLevel {
		return this.StandardFormatter.Format(entry)
	}

	content, err := this.StandardFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	// find the first caller that isn't in the logrus package
	var skip int
	for i, caller := range stack.Callers(1) {
		if !strings.Contains(caller.File, "/Sirupsen/logrus/") {
			skip = i + 1
			break
		}
	}

	trace := fmt.Sprintf("%s\n", stack.Callers(skip))
	content = append(content, []byte(trace)...)

	return content, nil
}

type SilentLogger struct{}

func (SilentLogger) Print(v ...interface{})                 {}
func (SilentLogger) Printf(format string, v ...interface{}) {}
