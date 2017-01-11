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

func (this *Reporter) Error(err error) {
	namedFormat := this.formatName(err.Error())
	this.T.Errorf(namedFormat)
}

func (this *Reporter) Errorf(format string, args ...interface{}) {
	namedFormat := this.formatName(format)
	this.T.Errorf(namedFormat, args...)
}

func (this *Reporter) Fatal(err error) {
	namedFormat := this.formatName(err.Error())
	this.T.Fatalf(namedFormat)
}

func (this *Reporter) Fatalf(format string, args ...interface{}) {
	namedFormat := this.formatName(format)
	this.T.Fatalf(namedFormat, args...)
}

func (this *Reporter) Log(str string) {
	namedFormat := this.formatName(str)
	this.T.Logf(namedFormat)
}

func (this *Reporter) Logf(format string, args ...interface{}) {
	namedFormat := this.formatName(format)
	this.T.Logf(namedFormat, args...)
}

func (this *Reporter) formatName(format string) string {
	return fmt.Sprintf("Test Case '%s': %s", this.Name, format)
}

func (this *Reporter) AssertEqual(result, expected interface{}) {
	this.AssertEqualf(result, expected, "")
}

func (this *Reporter) AssertAny(result interface{}, expected ...interface{}) {
	this.AssertInSlice(result, expected)
}

func (this *Reporter) AssertEqualf(result, expected interface{}, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)

	if !reflect.DeepEqual(expected, result) {
		message += fmt.Sprintf(
			"\n\tObserved: %#v (%v) \n\tExpected: %#v (%v)",
			result,
			reflect.TypeOf(result),
			expected,
			reflect.TypeOf(expected))

		this.Errorf(message)
	}
}

func (this *Reporter) AssertInSlice(expected, slice interface{}) {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		this.Errorf("\n\t%v (%v) is not a slice", slice, reflect.TypeOf(slice))
		return
	}

	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(expected, s.Index(i).Interface()) {
			return
		}
	}

	this.Errorf(
		"\n\tExpected Object: %#v (%v) \n\tIs not in slice: %#v (%v)",
		expected,
		reflect.TypeOf(expected),
		slice,
		reflect.TypeOf(slice))
}
