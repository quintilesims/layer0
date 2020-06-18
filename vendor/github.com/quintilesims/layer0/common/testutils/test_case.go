package testutils

import (
	"testing"

	"github.com/golang/mock/gomock"
)

type TestCase struct {
	Name  string
	Setup func(*Reporter, *gomock.Controller) interface{}
	Run   func(*Reporter, interface{})
}

func RunTest(t *testing.T, testCase TestCase) {
	reporter := NewReporter(t, testCase.Name)
	ctrl := gomock.NewController(reporter)
	defer ctrl.Finish()

	target := testCase.Setup(reporter, ctrl)
	testCase.Run(reporter, target)

	if !t.Failed() {
		t.Logf("PASS: %s", testCase.Name)
	}
}

func RunTests(t *testing.T, testCases []TestCase) {
	for _, testCase := range testCases {
		RunTest(t, testCase)
	}
}
