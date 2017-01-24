package system

import (
	"flag"
	"fmt"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/tests/system/framework"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

var (
	dry     = flag.Bool("dry", false, "Perform a dry run - don't execute terraform 'apply' commands")
	verbose = false
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

	sigtermChan := make(chan os.Signal)
	signal.Notify(sigtermChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigtermChan
		fmt.Println("\nInterrupt received. Shutting down")
		// todo: wait for subprocesses to die instead of just waiting here
		time.Sleep(time.Second * 1)
		os.Exit(1)
	}()
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
		T:       t,
		Dir:     dir,
		DryRun:  *dry,
		Verbose: verbose,
		Vars:    vars,
	})

	c.Apply()
	return c
}

func waitFor(t *testing.T, name string, timeout time.Duration, conditionSatisfied func() bool) {
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second * 1) {
		if verbose {
			fmt.Printf("Waiting for '%s' (waited %s of %s)\n", name, time.Since(start), timeout)
		}

		if conditionSatisfied() {
			return
		}
	}

	t.Fatalf("Wait for '%s' failed to complete after %v", name, timeout)
}
