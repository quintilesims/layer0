package main

import (
	"testing"

	"github.com/quintilesims/tftest"
	"github.com/stretchr/testify/assert"
)

func TestCaseTwo(t *testing.T) {
	t.Parallel()

	context := tftest.NewTestContext(t, tftest.Dir("./two"))
	context.Apply()
	defer context.Destroy()

	caseNumber := context.Output("case_number")
	assert.Equal(t, "two", caseNumber)
}
