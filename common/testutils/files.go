package testutils

import (
	"io/ioutil"
	"os"
	"testing"
)

func TempFile(t *testing.T, content string) (*os.File, func()) {
	file, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	delete := func() {
		if err := os.Remove(file.Name()); err != nil {
			t.Error(err)
		}
	}

	return file, delete
}
