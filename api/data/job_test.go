package data

import (
	"os"
	"testing"
)

func TestListJobs_successes(t *testing.T) {
}

func TestListJobs_errors(t *testing.T) {
	os.Unsetenv("LAYER0_SQLITE_DB_PATH")
}
