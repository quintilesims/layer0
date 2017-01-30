package system

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/logutils"
	"github.com/quintilesims/layer0/tests/system/framework"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logutils.SetGlobalLogger(log)
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

func startSystemTest(t *testing.T, dir string, vars map[string]string) *framework.SystemTestContext {
	t.Parallel()

	if vars == nil {
		vars = map[string]string{}
	}

	// add default terraform variables
	vars["endpoint"] = config.APIEndpoint()
	vars["token"] = config.AuthToken()

	c := framework.NewSystemTestContext(framework.Config{
		T:      t,
		Dir:    dir,
		DryRun: *dry,
		Vars:   vars,
	})

	c.Apply()
	return c
}

func waitFor(t *testing.T, name string, timeout time.Duration, conditionSatisfied func() bool) {
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second * 1) {
		log.Debugf("Waiting for '%s' (waited %s of %s)\n", name, time.Since(start), timeout)

		if conditionSatisfied() {
			return
		}
	}

	t.Fatalf("Wait for '%s' failed to complete after %v", name, timeout)
}
