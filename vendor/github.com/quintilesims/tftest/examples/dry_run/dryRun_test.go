package main

import (
	"os"
	"testing"

	"github.com/quintilesims/tftest"
)

func TestDryRun(t *testing.T) {
	if !*DryRun {
		t.Skip("-dry flag not set! Please run `go test -dry`")
	}

	context := tftest.NewTestContext(t, tftest.DryRun(*DryRun))
	context.Apply()
	defer context.Destroy()

	path := context.Output("file_path")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("File %s exists! Dry run failed", path)
	}
}
