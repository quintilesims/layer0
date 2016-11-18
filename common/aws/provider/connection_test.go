package provider

import (
	"errors"
	"testing"
)

func TestInvalidRegion(t *testing.T) {
	creds := MockCredProvider{}
	_, err := getConfig(&creds, "invalid_region")
	if err == nil {
		t.Error("No error returned!")
	}
}

func TestCredProviderErrorPropagates(t *testing.T) {
	creds := MockCredProvider{}
	creds.GetAWSSecretAccessKey_fn = func() (string, error) {
		return "", errors.New("some error")
	}
	_, err := getConfig(&creds, US_WEST_1)
	if err == nil {
		t.Error("No error returned!")
	}
}
