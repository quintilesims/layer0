package system

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/logutils"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	dry   = flag.Bool("dry", false, "Perform a dry run - don't execute terraform 'apply' commands")
	debug = flag.Bool("debug", false, "Print debug statements")
	log   = logutils.NewStandardLogger("Test")
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	flag.Parse()

	log.Level = logrus.ErrorLevel
	if *debug {
		log.Level = logrus.DebugLevel
	}

	logutils.SetGlobalLogger(log)

	if !*dry {
		if err := filepath.Walk("cases", deleteStateFiles); err != nil {
			fmt.Println("Error occurred during setup: ", err)
			os.Exit(1)
		}
	}
}

func teardown() {
	if !*dry {
		if err := filepath.Walk("cases", deleteStateFiles); err != nil {
			fmt.Println("Error occurred during teardown: ", err)
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
			return fmt.Errorf("Failed to delete %s: %v", path, err)
		}
	}

	return nil
}
