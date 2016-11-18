package testutils

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, result, expected interface{}) {
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(
			"\n\tObserved: %#v (%v) \n\tExpected: %#v (%v)",
			result,
			reflect.TypeOf(result),
			expected,
			reflect.TypeOf(expected))
	}
}

func AssertInSlice(t *testing.T, expected, slice interface{}) {
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
