package context

import (
	"fmt"
	"os"
	"testing"
)

type MockFileIO struct {
	data map[string][]byte
}

func (mock MockFileIO) ReadFile(path string) ([]byte, error) {
	if contents, ok := mock.data[path]; ok {
		return contents, nil
	}
	return nil, fmt.Errorf("Failed to read file at path: `%s`", path)
}

func (mock MockFileIO) WriteFile(path string, data []byte, perm os.FileMode) error {
	mock.data[path] = data
	return nil
}

func (mock MockFileIO) Stat(path string) (os.FileInfo, error) {
	if _, ok := mock.data[path]; ok {
		return nil, nil
	}
	return nil, fmt.Errorf("Failed to stat file at path: `%s`", path)
}

func TestValidateDockercfgWithValidData(t *testing.T) {
	mockIO := MockFileIO{data: map[string][]byte{"foo": []byte(`
        {
            "https://d.ims.io": {
                "auth": "StopLookingAtMySecretsYouJerk=",
                "email": ""
            }
        }`)}}
	if err := validateDockercfg("foo", mockIO); err != nil {
		t.Fatal(err)
	}
}

func TestValidateDockercfgWithInvalidData(t *testing.T) {
	mockIO := MockFileIO{data: map[string][]byte{"foo": []byte(`
        {
            "https://d.ims.io": {
                "auth": "StopLookingAtMySecretsYouJerk=",
                "email": ""
            },
        }`)}}
	if err := validateDockercfg("foo", mockIO); err == nil {
		t.Fatal(err)
	}
}
