package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	cases := []struct {
		Name           string
		Error          error
		ExpectedStatus int
		ExpectedCode   errors.ErrorCode
	}{
		{
			Name:           "Basic string error",
			Error:          fmt.Errorf("some error"),
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedCode:   errors.UnexpectedError,
		},
		{
			Name:           "Unexpected ServerError",
			Error:          errors.Newf(errors.UnexpectedError, ""),
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedCode:   errors.UnexpectedError,
		},
		{
			Name:           "ServerError: InvalidRequest",
			Error:          errors.Newf(errors.InvalidRequest, ""),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedCode:   errors.InvalidRequest,
		},
		{
			Name:           "ServerError: DeployDoesNotExist",
			Error:          errors.Newf(errors.DeployDoesNotExist, ""),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedCode:   errors.DeployDoesNotExist,
		},
		{
			Name:           "ServerError: EnvironmentDoesNotExist",
			Error:          errors.Newf(errors.EnvironmentDoesNotExist, ""),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedCode:   errors.EnvironmentDoesNotExist,
		},
		{
			Name:           "ServerError: JobDoesNotExist",
			Error:          errors.Newf(errors.JobDoesNotExist, ""),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedCode:   errors.JobDoesNotExist,
		},
		{
			Name:           "ServerError: LoadBalancerDoesNotExist",
			Error:          errors.Newf(errors.LoadBalancerDoesNotExist, ""),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedCode:   errors.LoadBalancerDoesNotExist,
		},
		{
			Name:           "ServerError: TaskDoesNotExist",
			Error:          errors.Newf(errors.TaskDoesNotExist, ""),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedCode:   errors.TaskDoesNotExist,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("", "", nil)
			if err != nil {
				t.Fatal(err)
			}

			ErrorHandler(recorder, req, c.Error)

			var result models.ServerError
			if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.ExpectedStatus, recorder.Code)
			assert.Equal(t, c.ExpectedCode.String(), result.ErrorCode)
		})
	}
}
