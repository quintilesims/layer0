package testutils

import (
	"reflect"
	"time"
)

func AssertEqual(t Tester, result, expected interface{}) {
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(
			"\n\tObserved: %#v (%v) \n\tExpected: %#v (%v)",
			result,
			reflect.TypeOf(result),
			expected,
			reflect.TypeOf(expected))
	}
}

func AssertAny(t Tester, result interface{}, expected ...interface{}) {
	AssertInSlice(t, result, expected)
}

func AssertInSlice(t Tester, expected, slice interface{}) {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		t.Fatalf("\n\t%v (%v) is not a slice", slice, reflect.TypeOf(slice))
		return
	}

	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(expected, s.Index(i).Interface()) {
			return
		}
	}

	t.Fatalf(
		"\n\tExpected Object: %#v (%v) \n\tIs not in slice: %#v (%v)",
		expected,
		reflect.TypeOf(expected),
		slice,
		reflect.TypeOf(slice))
}

func WaitFor(t Tester, interval, timeout time.Duration, conditionSatisfied func() bool) {
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(interval) {
		if conditionSatisfied() {
			return
		}
	}

	t.Fatalf("Timout reached after %v", timeout)
}
