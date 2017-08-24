package job

type JobType string

const (
	CreateDeployJob       JobType = "CreateDeploy"
	CreateEnvironmentJob  JobType = "CreateEnvironment"
	CreateLoadBalancerJob JobType = "CreateLoadBalancer"
	CreateServiceJob      JobType = "CreateService"
	CreateTaskJob         JobType = "CreateTask"
	DeleteDeployJob       JobType = "DeleteDeploy"
	DeleteEnvironmentJob  JobType = "DeleteEnvironment"
	DeleteLoadBalancerJob JobType = "DeleteLoadBalancer"
	DeleteServiceJob      JobType = "DeleteService"
	DeleteTaskJob         JobType = "DeleteTask"
)

type JobStatus string

const (
	Pending    JobStatus = "Pending"
	InProgress JobStatus = "InProgress"
	Completed  JobStatus = "Completed"
	Error      JobStatus = "Error"
)
