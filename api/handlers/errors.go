package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func ToHttpError(code errors.ErrorCode) int {
	var ret int
	switch code {
	case errors.InvalidJSON, errors.MissingParameter, errors.InvalidEntityType,
		errors.InvalidEnvironmentID, errors.InvalidServiceID, errors.InvalidDeployID,
		errors.InvalidTagKey, errors.InvalidTagValue, errors.InvalidCertificateID:
		ret = http.StatusBadRequest
	case errors.Throttled:
		ret = http.StatusServiceUnavailable
	default:
		ret = http.StatusInternalServerError
	}

	return ret
}

func ReturnError(response *restful.Response, err error) {
	switch err := err.(type) {
	case awserr.Error:
		AWSError(response, err)
	case *errors.ServerError:
		Errorf(response, ToHttpError(err.Code), err.Model())
	default:
		internalServerError(response, errors.UnexpectedError, err)
	}
}

func BadRequest(response *restful.Response, code errors.ErrorCode, err error) {
	logrus.Infof("BadRequest: %s (code: %d)", err.Error(), code)
	model := models.ServerError{ErrorCode: int64(code), Message: err.Error()}
	Errorf(response, http.StatusBadRequest, model)
}

func AWSError(response *restful.Response, err awserr.Error) {
	logrus.Errorf("AWSError (code: %s): %s", err.Code(), err.Message())

	var code errors.ErrorCode
	var status int

	switch err.Code() {
	case "Throttling":
		code = errors.Throttled
		status = http.StatusServiceUnavailable
	default:
		code = errors.UnexpectedError
		status = http.StatusInternalServerError
	}

	message := strings.ToLower(err.Message())
	message = strings.Replace(message, "environment", "service", -1)
	message = strings.Replace(message, "application", "environment", -1)
	message = strings.Replace(message, "listener", "port", -1)

	message = fmt.Sprintf("AWS Error: %s (code '%s')", message, err.Code())
	model := models.ServerError{ErrorCode: int64(code), Message: message}
	Errorf(response, status, model)
}

func NotImplemented(request *restful.Request, response *restful.Response) {
	logrus.Errorf("NotImplementedError")
	model := models.ServerError{ErrorCode: int64(501), Message: "Not Implemented"}
	Errorf(response, http.StatusNotImplemented, model)
}

func NotFound(response *restful.Response, message string) {
	model := models.ServerError{ErrorCode: int64(404), Message: message}
	Errorf(response, http.StatusNotFound, model)
}

func Errorf(response *restful.Response, status int, se models.ServerError) {
	response.WriteHeader(status)
	response.WriteAsJson(se)
}

func internalServerError(response *restful.Response, code errors.ErrorCode, err error) {
	logrus.Errorf("InternalServerError (code: %d): %s", code, err.Error())
	model := models.ServerError{ErrorCode: int64(code), Message: err.Error()}
	Errorf(response, http.StatusInternalServerError, model)
}
