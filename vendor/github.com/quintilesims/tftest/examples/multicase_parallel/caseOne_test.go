package main

import (
	"testing"

	"github.com/quintilesims/tftest"
	"github.com/stretchr/testify/assert"
)

func TestCaseOne(t *testing.T) {
	t.Parallel()

	context := tftest.NewTestContext(t, tftest.Dir("./one"))
	context.Apply()
	defer context.Destroy()

	caseNumber := context.Output("case_number")
	assert.Equal(t, "one", caseNumber)
}
