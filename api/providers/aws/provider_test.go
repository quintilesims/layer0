package aws

import (
	"testing"

	"github.com/quintilesims/layer0/common/errors"
	"github.com/stretchr/testify/assert"
)

func assertEqualErrorCode(t *testing.T, err error, code errors.ErrorCode) {
	switch err := err.(type) {
	case *errors.ServerError:
		assert.Equal(t, code, err.Code)
	default:
		t.Fatal("Error %#v is not of type *errors.ServerError", err)
	}
}
