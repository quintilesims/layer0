package system

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	dry     = flag.Bool("dry", false, "Perform a dry run - don't execute terraform 'apply' commands")
	verbose bool
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	flag.Parse()

	// use `go test -v` flag to determine verbosity
	if v := flag.Lookup("test.v"); v != nil {
		verbose = v.Value.String() == "true"
	}
}

func teardown() {
	deleteStateFiles := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if name := f.Name(); strings.HasPrefix(name, "terraform.tfstate") {
			if err := os.Remove(path); err != nil {
				fmt.Println("Error occurred during teardown: ", err)
			}
		}

		return nil
	}

	if err := filepath.Walk("cases", deleteStateFiles); err != nil {
		fmt.Println("Error occurred during teardown: ", err)
	}
}
