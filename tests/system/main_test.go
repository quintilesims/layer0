package system

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func runSystemTest(t *testing.T, dir string) *SystemTestContext {
	t.Parallel()
	c := NewSystemTestContext(t, "example")
	c.Apply()
	return c
}
