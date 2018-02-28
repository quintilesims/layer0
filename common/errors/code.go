package errors

type ErrorCode string

func (e ErrorCode) String() string {
	return string(e)
}

const (
	InvalidRequest                ErrorCode = "InvalidReqest"
	DeployDoesNotExist            ErrorCode = "DeployDoesNotExist"
	EnvironmentDoesNotExist       ErrorCode = "EnvironmentDoesNotExist"
	DependencyError               ErrorCode = "DependencyError"
	JobDoesNotExist               ErrorCode = "JobDoesNotExist"
	LoadBalancerDoesNotExist      ErrorCode = "LoadBalancerDoesNotExist"
	LoadBalancerAttributeNotFound ErrorCode = "LoadBalancerAttributeNotFound"
	ServiceDoesNotExist           ErrorCode = "ServiceDoesNotExist"
	TaskDoesNotExist              ErrorCode = "TaskDoesNotExist"
	UnexpectedError               ErrorCode = "UnexpectedError"
)
