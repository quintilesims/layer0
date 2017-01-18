package system

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// todo: setup
	os.Exit(m.Run())	
	// todo: teardown - remove all .tfstate* files
}

func runSystemTest(t *testing.T, dir string) *SystemTestContext {
	t.Parallel()
	c := NewSystemTestContext(t, dir)
	c.Apply()
	return c
}
