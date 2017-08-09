package main

import (
	"testing"

	"github.com/quintilesims/tftest"
)

func TestHelloWorld(t *testing.T) {
	context := tftest.NewTestContext(t)
	context.Apply()
	defer context.Destroy()

	message := context.Output("message")
	if message != "Hello World" {
		t.Errorf("message was '%s', expected 'Hello World'", message)
	}
}
