package errors

import (
	"fmt"
)

func MultiError(errors []error) error {
	if count := len(errors); count == 0 {
		return nil
	} else if count == 1 {
		return errors[0]
	}

	text := "Multiple Errors: \n"
	for _, err := range errors {
		text += fmt.Sprintf("\t%s\n", err.Error())
	}

	return fmt.Errorf(text)
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
