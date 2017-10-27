package main

import (
	"testing"

	"github.com/quintilesims/tftest"
	"github.com/stretchr/testify/assert"
)

func TestCaseOne(t *testing.T) {
	context := tftest.NewTestContext(t, tftest.Vars(map[string]string{
		"message_one": "Hello",
		"message_two": "World",
	}))

	context.Apply()
	defer context.Destroy()

	assert.Equal(t, "Hello World", context.Output("combined_messages"))
}
