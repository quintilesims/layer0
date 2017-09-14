package errors

type ErrorCode string

func (e ErrorCode) String() string {
	return string(e)
}

const (
	InvalidRequest           ErrorCode = "InvalidReqest"
	DeployDoesNotExist       ErrorCode = "DeployDoesNotExist"
	EnvironmentDoesNotExist  ErrorCode = "EnvironmentDoesNotExist"
	JobDoesNotExist          ErrorCode = "JobDoesNotExist"
	LoadBalancerDoesNotExist           = "LoadBalancerDoesNotExist"
	ServiceDoesNotExist                = "ServiceDoesNotExist"
	TaskDoesNotExist                   = "TaskDoesNotExist"
	UnexpectedError                    = "UnexpectedError"
)
