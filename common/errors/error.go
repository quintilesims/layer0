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

func NewEntityDoesNotExistError(entityType, entityID string) *ServerError {
	switch entityType {
	case "deploy":
		return Newf(DeployDoesNotExist, "Deploy '%s' does not exist", entityID)
	case "environment":
		return Newf(EnvironmentDoesNotExist, "Environment '%s' does not exist", entityID)
	case "job":
		return Newf(JobDoesNotExist, "Job '%s' does not exist", entityID)
	case "load_balancer":
		return Newf(LoadBalancerDoesNotExist, "Load balancer '%s' does not exist", entityID)
	case "service":
		return Newf(ServiceDoesNotExist, "Service '%s' does not exist", entityID)
	case "task":
		return Newf(TaskDoesNotExist, "Task '%s' does not exist", entityID)
	default:
		return Newf(UnexpectedError, "Entity (type='%s') '%s' does not exist", entityType, entityID)
	}
}
