package context

import (
	"os"
	"testing"
)

type MockFileIO struct {
	data []byte
}

func (mock MockFileIO) ReadFile(path string) ([]byte, error) {
	return mock.data, nil
}

func (mock MockFileIO) WriteFile(path string, data []byte, perm os.FileMode) error {
	mock.data = data
	return nil
}

func (mock MockFileIO) Stat(path string) (os.FileInfo, error) {
	return nil, nil
}

func TestValidateDockercfgWithValidData(t *testing.T) {
	mockIO := MockFileIO{data: []byte(`
        {
            "https://d.ims.io": {
                "auth": "StopLookingAtMySecretsYouJerk=",
                "email": ""
            }
        }`)}
	if err := validateDockercfg("foo", mockIO); err != nil {
		t.Fatal(err)
	}
}

func TestValidateDockercfgWithInvalidData(t *testing.T) {
	mockIO := MockFileIO{data: []byte(`
        {
            "https://d.ims.io": {
                "auth": "StopLookingAtMySecretsYouJerk=",
                "email": ""
            },
        }`)}
	if err := validateDockercfg("foo", mockIO); err == nil {
		t.Fatal(err)
	}
}
