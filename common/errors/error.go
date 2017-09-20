package errors

import (
	"fmt"

	"github.com/quintilesims/layer0/common/models"
)

type ServerError struct {
	Code ErrorCode
	Err  error
}

func Newf(code ErrorCode, format string, tokens ...interface{}) *ServerError {
	return New(code, fmt.Errorf(format, tokens...))
}

func New(code ErrorCode, err error) *ServerError {
	return &ServerError{
		Code: code,
		Err:  err,
	}
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("ServerError (code='%s') %s", s.Code, s.Err.Error())
}

func (s *ServerError) Model() models.ServerError {
	return models.ServerError{
		ErrorCode: s.Code.String(),
		Message:   s.Err.Error(),
	}
}

func ResolveErrorByEntityType(entityType string) (ErrorCode, error) {
	switch entityType {
	case "deploy":
		return DeployDoesNotExist, nil
	case "environment":
		return EnvironmentDoesNotExist, nil
	case "job":
		return JobDoesNotExist, nil
	case "load_balancer":
		return LoadBalancerDoesNotExist, nil
	case "service":
		return ServiceDoesNotExist, nil
	case "task":
		return TaskDoesNotExist, nil
	default:
		return "", Newf(InvalidRequest, "Entity type '%s' is invalid", entityType)
	}
}
