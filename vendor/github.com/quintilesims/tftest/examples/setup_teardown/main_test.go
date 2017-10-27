package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/quintilesims/tftest"
)

func TestMain(m *testing.M) {
	context := tftest.NewContext(tftest.Dir("./setup"))

	setup(context)
	code := m.Run()
	teardown(context)

	os.Exit(code)
}

func setup(context *tftest.Context) {
	fmt.Println("Running Setup")
	if _, err := context.Apply(); err != nil {
		fmt.Printf("Error during setup: %v", err)
		os.Exit(1)
	}
}

func teardown(context *tftest.Context) {
	fmt.Println("Running Teardown")
	if _, err := context.Destroy(); err != nil {
		fmt.Printf("Error during teardown: %v", err)
		os.Exit(1)
	}
}
