package errors

type ErrorCode int64

const (
	InvalidJSON ErrorCode = 1 + iota
	InvalidCertificateID
	InvalidCertificateName
	InvalidDeployID
	InvalidDeployName
	InvalidEnvironmentID
	InvalidEnvironmentName
	InvalidJobID
	InvalidJobName
	InvalidLoadBalancerID
	InvalidLoadBalancerName
	InvalidServiceID
	InvalidServiceName
	InvalidTaskID
	InvalidTaskName
	InvalidEntityType
	InvalidTagKey
	InvalidTagValue
	MissingParameter
	Throttled
	UnexpectedError
	FailedTagging
)
