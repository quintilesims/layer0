package tftest

type Logger interface {
	Printf(string, ...interface{})
}

type TestLogger struct {
	t Tester
}

func NewTestLogger(t Tester) *TestLogger {
	return &TestLogger{t: t}
}

func (l *TestLogger) Printf(format string, tokens ...interface{}) {
	l.t.Logf(format, tokens...)
}
