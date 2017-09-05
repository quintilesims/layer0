package testutils

type Tester interface {
	Fatal(tokens ...interface{})
	Fatalf(format string, tokens ...interface{})
}
