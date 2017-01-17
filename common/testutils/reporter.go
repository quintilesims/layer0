package testutils

import (
	"fmt"
	"reflect"
	"testing"
)

type Reporter struct {
	T    *testing.T
	Name string
}

func NewReporter(t *testing.T, name string) *Reporter {
	return &Reporter{
		T:    t,
		Name: name,
	}
}

func (r *Reporter) Error(err error) {
	namedFormat := r.formatName(err.Error())
	r.T.Errorf(namedFormat)
}

func (r *Reporter) Errorf(format string, args ...interface{}) {
	namedFormat := r.formatName(format)
	r.T.Errorf(namedFormat, args...)
}

func (r *Reporter) Fatal(err error) {
	namedFormat := r.formatName(err.Error())
	r.T.Fatalf(namedFormat)
}

func (r *Reporter) Fatalf(format string, args ...interface{}) {
	namedFormat := r.formatName(format)
	r.T.Fatalf(namedFormat, args...)
}

func (r *Reporter) Log(str string) {
	namedFormat := r.formatName(str)
	r.T.Logf(namedFormat)
}

func (r *Reporter) Logf(format string, args ...interface{}) {
	namedFormat := r.formatName(format)
	r.T.Logf(namedFormat, args...)
}

func (r *Reporter) formatName(format string) string {
	return fmt.Sprintf("Test Case '%s': %s", r.Name, format)
}

func (r *Reporter) AssertEqual(result, expected interface{}) {
	r.AssertEqualf(result, expected, "")
}

func (r *Reporter) AssertAny(result interface{}, expected ...interface{}) {
	r.AssertInSlice(result, expected)
}

func (r *Reporter) AssertEqualf(result, expected interface{}, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)

	if !reflect.DeepEqual(expected, result) {
		message += fmt.Sprintf(
			"\n\tObserved: %#v (%v) \n\tExpected: %#v (%v)",
			result,
			reflect.TypeOf(result),
			expected,
			reflect.TypeOf(expected))

		r.Fatalf(message)
	}
}

func (r *Reporter) AssertInSlice(expected, slice interface{}) {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		r.Errorf("\n\t%v (%v) is not a slice", slice, reflect.TypeOf(slice))
		return
	}

	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(expected, s.Index(i).Interface()) {
			return
		}
	}

	r.Fatalf(
		"\n\tExpected Object: %#v (%v) \n\tIs not in slice: %#v (%v)",
		expected,
		reflect.TypeOf(expected),
		slice,
		reflect.TypeOf(slice))
}
