package errors

type ErrorCode int64

const (
	InvalidRequest ErrorCode = 1 + iota
	DeployDoesNotExist
	EnvironmentDoesNotExist
	JobDoesNotExist
	LoadBalancerDoesNotExist
	ServiceDoesNotExist
	TaskDoesNotExist
	UnexpectedError
)
