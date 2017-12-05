package system

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/quintilesims/layer0/common/logging"
)

var (
	dry    = flag.Bool("dry", false, "Perform a dry run - don't execute terraform 'apply' commands")
	debug  = flag.Bool("debug", false, "Print debug statements")
	logger = logging.NewLogWriter(*debug)
)

func TestMain(m *testing.M) {
	log.SetOutput(logger)
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	flag.Parse()
	if !*dry {
		if err := filepath.Walk("cases", deleteStateFiles); err != nil {
			log.Fatalf("[ERROR] Error occurred during setup: ", err)
			os.Exit(1)
		}
	}
}

func teardown() {
	if !*dry {
		if err := filepath.Walk("cases", deleteStateFiles); err != nil {
			log.Fatalf("[ERROR] Error occurred during teardown: ", err)
			os.Exit(1)
		}
	}
}

func deleteStateFiles(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if name := f.Name(); strings.HasPrefix(name, "terraform.tfstate") {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("[ERROR] Failed to delete %s: %v", path, err)
		}
	}

	return nil
}
