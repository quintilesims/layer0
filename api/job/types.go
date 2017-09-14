package job

type JobType string

func (j JobType) String() string {
	return string(j)
}

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
	UpdateEnvironmentJob  JobType = "UpdateEnvironment"
)

type Status string

func (j Status) String() string {
	return string(j)
}

const (
	Pending    Status = "Pending"
	InProgress Status = "InProgress"
	Completed  Status = "Completed"
	Error      Status = "Error"
)
